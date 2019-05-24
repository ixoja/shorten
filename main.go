package main

import (
	"log"
	"net"
	"net/http"

	"github.com/ixoja/shorten/internal/grpcapi"
	"github.com/ixoja/shorten/internal/service"
	"google.golang.org/grpc"

	"github.com/ixoja/shorten/internal/webserver"
)

func main() {
	const webServer = "webserver"
	const grpcServer = "grpcserver"
	config := Config{}
	config.WithFlags()

	switch config.mode {
	case webServer:
		conn, err := grpc.Dial(config.apiURL, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer func() {
			if err := conn.Close(); err != nil {
				log.Println("failed to start grpc client:", err.Error())
			}
		}()

		c := grpcapi.NewShortenServiceClient(conn)
		ws := webserver.New(c, config.webURL)
		http.Handle("/", http.FileServer(http.Dir(config.htmlPath)))
		http.HandleFunc("/shorten", ws.Shorten)
		http.HandleFunc("/to", ws.Redirect)
		log.Println("Registering web server on port:", config.port)
		log.Fatal(http.ListenAndServe(":"+config.port, nil))
	case grpcServer:
		lis, err := net.Listen("tcp", config.port)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		s := grpc.NewServer()
		grpcapi.RegisterShortenServiceServer(s, &service.Service{})
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}
}
