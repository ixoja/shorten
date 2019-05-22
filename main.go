package main

import (
	"github.com/ixoja/shorten/internal/grpcapi"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"

	"github.com/ixoja/shorten/internal/webserver"
)

func main() {
	const webServer = "webserver"
	const grpcServer = "grpcserver"
	config := Config{}
	config.WithFlags()

	switch config.mode {
	case webServer:
		c := webserver.NewClient(http.DefaultClient)
		ws := webserver.New(c, config.apiURL)
		http.Handle("/", http.FileServer(http.Dir(config.htmlPath)))
		http.HandleFunc("/shorten", ws.Shorten)
		log.Println("Registering web server on port:", config.port)
		log.Fatal(http.ListenAndServe(":"+config.port, nil))
	case grpcServer:
		lis, err := net.Listen("tcp", config.port)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		s := grpc.NewServer()
		grpcapi.RegisterShortenServiceServer(s, &server{})
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}
}
