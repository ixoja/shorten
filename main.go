package main

import (
	"database/sql"
	"github.com/ixoja/shorten/internal/controller"
	"github.com/ixoja/shorten/internal/storage"
	"log"
	"net"
	"net/http"

	"github.com/ixoja/shorten/internal/grpcapi"
	"github.com/ixoja/shorten/internal/service"
	_ "github.com/mattn/go-sqlite3"
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
		db, err := sql.Open("sqlite3", "./foo.db")
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			if err := db.Close(); err != nil {
				log.Println("failed to start db connection:", err.Error())
			}
		}()

		s := service.New(controller.New(&storage.Cache{}, storage.New(*db)))
		lis, err := net.Listen("tcp", config.port)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		server := grpc.NewServer()
		grpcapi.RegisterShortenServiceServer(server, s)
		if err := server.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}
}
