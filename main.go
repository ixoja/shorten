package main

import (
	"flag"
	"github.com/ixoja/shorten/internal/handler"
	"github.com/ixoja/shorten/internal/restapi"
	"github.com/ixoja/shorten/internal/webserver"
	"log"
	"os"
)

func main() {
	const webServer = "webserver"
	const server = "server"
	mode := flag.String("mode", webServer, "string value: webapi or server")
	flag.Parse()

	switch *mode {
	case webServer:
		htmlPath := os.Getenv("HTML")
		if htmlPath == "" {
			log.Fatal("HTML value not set")
		}
		port := os.Getenv("PORT")
		if port == "" {
			log.Fatal("PORT value not set")
		}
		ws := webserver.New(htmlPath, port)
		ws.Start()
	case server:
		var apiConfig restapi.Config
		err := apiConfig.Parse()
		if err != nil {
			log.Fatal(err)
		}
		svc := &handler.Service{}
		api := restapi.NewServer(svc, &apiConfig)
		err = api.RunWithSigHandler()
		if err != nil {
			log.Fatal(err)
		}
	}
}
