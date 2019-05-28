package main

import (
	"database/sql"
	"log"
	"net"
	"net/http"

	"github.com/ixoja/shorten/internal/controller"
	"github.com/ixoja/shorten/internal/grpcapi"
	"github.com/ixoja/shorten/internal/service"
	"github.com/ixoja/shorten/internal/storage"
	"github.com/ixoja/shorten/internal/webserver"
	_ "github.com/mattn/go-sqlite3"
	"google.golang.org/grpc"
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
		ws := webserver.New(c, config.webURL+":"+config.port)
		http.Handle("/", http.FileServer(http.Dir(config.htmlPath)))
		http.HandleFunc("/shorten", ws.Shorten)
		http.HandleFunc("/to", ws.Redirect)
		log.Println("Registering web server on port:", config.port)
		log.Fatal(http.ListenAndServe(":"+config.port, nil))
	case grpcServer:
		db, err := sql.Open("sqlite3", "./shorten.db")
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			if err := db.Close(); err != nil {
				log.Println("failed to start db connection:", err.Error())
			}
		}()

		st := storage.New(db)
		if err := st.InitDB(); err != nil {
			log.Fatalf("failed to init db: %v", err)
		}

		s := service.New(controller.New(storage.NewCache(), st))
		lis, err := net.Listen("tcp", config.port)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		server := grpc.NewServer()
		grpcapi.RegisterShortenServiceServer(server, s)
		log.Println("strating grpc service")
		if err := server.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}
}
