package main

import (
	"github.com/ixoja/shorten/internal/handler"
	"github.com/ixoja/shorten/internal/restapi"
	"github.com/ixoja/shorten/internal/webserver"
	"log"
	"net/http"
)

func main() {
	const webServer = "webserver"
	const apiServer = "apiserver"
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
	case apiServer:
		svc := &handler.Service{}
		api := restapi.NewServer(svc, &config.apiConfig)
		err := api.RunWithSigHandler()
		if err != nil {
			log.Fatal(err)
		}
	}
}
