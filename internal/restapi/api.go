package restapi

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/ixoja/shorten/internal/restapi/operations"
	"github.com/mikkeloscar/gin-swagger/api"
	"github.com/mikkeloscar/gin-swagger/middleware"
	"github.com/mikkeloscar/gin-swagger/tracing"
	opentracing "github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
)

// Routes defines all the routes of the Server service.
type Routes struct {
	*gin.Engine
	GetHash struct {
		*gin.RouterGroup
	}
	PostShorten struct {
		*gin.RouterGroup
	}
}

// configureWellKnown enables and configures /.well-known endpoints.
func (r *Routes) configureWellKnown(healthFunc func() bool) {
	wellKnown := r.Group("/.well-known")
	{
		wellKnown.GET("/schema-discovery", func(ctx *gin.Context) {
			discovery := struct {
				SchemaURL  string `json:"schema_url"`
				SchemaType string `json:"schema_type"`
				UIURL      string `json:"ui_url"`
			}{
				SchemaURL:  "/swagger.json",
				SchemaType: "swagger-2.0",
			}
			ctx.JSON(http.StatusOK, &discovery)
		})
		wellKnown.GET("/health", healthHandler(healthFunc))
	}

	r.GET("/swagger.json", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, string(SwaggerJSON))
	})
}

// healthHandler is the health HTTP handler used for the /.well-known/health
// route if enabled.
func healthHandler(healthFunc func() bool) gin.HandlerFunc {
	healthy := healthFunc
	return func(ctx *gin.Context) {
		health := struct {
			Health bool `json:"health"`
		}{
			Health: healthy(),
		}

		if !health.Health {
			ctx.JSON(http.StatusServiceUnavailable, &health)
		} else {
			ctx.JSON(http.StatusOK, &health)
		}
	}
}

// Service is the interface that must be implemented in order to provide
// business logic for the Server service.
type Service interface {
	Healthy() bool
	GetHash(ctx *gin.Context, params *operations.GetHashParams) *api.Response
	PostShorten(ctx *gin.Context, params *operations.PostShortenParams) *api.Response
}

func ginizePath(path string) string {
	return strings.Replace(strings.Replace(path, "{", ":", -1), "}", "", -1)
}

// initializeRoutes initializes the route structure for the Server service.
func initializeRoutes(enableAuth bool, tokenURL string, tracer opentracing.Tracer) *Routes {
	engine := gin.New()
	engine.Use(gin.Recovery())
	routes := &Routes{Engine: engine}

	routes.GetHash.RouterGroup = routes.Group("")
	routes.GetHash.RouterGroup.Use(middleware.LogrusLogger())
	if tracer != nil {
		routes.GetHash.RouterGroup.Use(tracing.InitSpan(tracer, "get_hash"))
	}

	routes.PostShorten.RouterGroup = routes.Group("")
	routes.PostShorten.RouterGroup.Use(middleware.LogrusLogger())
	if tracer != nil {
		routes.PostShorten.RouterGroup.Use(tracing.InitSpan(tracer, "post_shorten"))
	}
	routes.PostShorten.RouterGroup.Use(middleware.ContentTypes("text/plain"))

	return routes
}

// Server defines the Server service.
type Server struct {
	Routes           *Routes
	config           *Config
	server           *http.Server
	service          Service
	healthy          bool
	serviceHealthyFn func() bool
	authDisabled     bool
	Title            string
	Version          string
}

// NewServer initializes a new Server service.
func NewServer(svc Service, config *Config) *Server {
	if !config.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	server := &Server{
		Routes: initializeRoutes(
			!config.AuthDisabled,
			config.TokenURL,
			config.Tracer,
		),
		service:      svc,
		config:       config,
		Title:        "URL Shortener API",
		Version:      "1.0",
		authDisabled: config.AuthDisabled,
	}

	// enable pprof http endpoints in debug mode
	if config.Debug {
		pprof.Register(server.Routes.Engine)
	}

	// set logrus logger to TextFormatter with no colors
	log.SetFormatter(&log.TextFormatter{DisableColors: true})

	server.server = &http.Server{
		Addr:         config.Address,
		Handler:      server.Routes.Engine,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	server.serviceHealthyFn = svc.Healthy

	if !config.WellKnownDisabled {
		server.Routes.configureWellKnown(server.isHealthy)
	}

	// configure healthz endpoint
	server.Routes.GET("/healthz", healthHandler(server.isHealthy))

	return server
}

// isHealthy returns true if both the server and the service reports healthy.
func (s *Server) isHealthy() bool {
	return s.healthy && s.serviceHealthyFn()
}

// ConfigureRoutes starts the internal configureRoutes methode.
func (s *Server) ConfigureRoutes() {
	s.configureRoutes()
}

// configureRoutes configures the routes for the Server service.
// Configuring of routes includes setting up Auth if it is enabled.
func (s *Server) configureRoutes() {
	if !s.authDisabled {
	}

	// setup all service routes after the authenticate middleware has been
	// initialized.
	s.Routes.GetHash.GET(ginizePath("/{hash}"), operations.GetHashEndpoint(s.service.GetHash))
	s.Routes.PostShorten.POST(ginizePath("/shorten"), operations.PostShortenEndpoint(s.service.PostShorten))
}

// Run runs the Server. It will listen on either HTTP or HTTPS depending on the
// config passed to NewServer.
func (s *Server) Run() error {
	// configure service routes
	s.configureRoutes()

	log.Infof("Serving '%s - %s' on address %s", s.Title, s.Version, s.server.Addr)
	// server is set to healthy when started.
	s.healthy = true
	if s.config.InsecureHTTP {
		return s.server.ListenAndServe()
	}
	return s.server.ListenAndServeTLS(s.config.TLSCertFile, s.config.TLSKeyFile)
}

// Shutdown will gracefully shutdown the Server server.
func (s *Server) Shutdown() error {
	// server is set to unhealthy when shutting down
	s.healthy = false
	return s.server.Shutdown(context.Background())
}

// RunWithSigHandler runs the Server server with SIGTERM handling automatically
// enabled. The server will listen for a SIGTERM signal and gracefully shutdown
// the web server.
// It's possible to optionally pass any number shutdown functions which will
// execute one by one after the webserver has been shutdown successfully.
func (s *Server) RunWithSigHandler(shutdown ...func() error) error {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigCh
		s.Shutdown()
	}()

	err := s.Run()
	if err != nil {
		if err != http.ErrServerClosed {
			return err
		}
	}

	for _, fn := range shutdown {
		err := fn()
		if err != nil {
			return err
		}
	}

	return nil
}

// vim: ft=go
