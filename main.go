package main

import (
	"flag"
	"github.com/ixoja/shorten/internal/handler"
	"github.com/ixoja/shorten/internal/restapi"
	"github.com/ixoja/shorten/internal/webserver"
	"log"
	"net/http"
)

func main() {
	const webServer = "webserver"
	const server = "server"
	mode := flag.String("mode", webServer, "string value: webserver or server")
	htmlPath := flag.String("html", "", "path to html dir")
	port := flag.String("port", "", "webserver http port")
	apiURL := flag.String("api_url", "", "web api url")
	flag.Parse()

	switch *mode {
	case webServer:

		c := webserver.NewClient(http.DefaultClient)
		ws := webserver.New(c, *apiURL)
		http.Handle("/", http.FileServer(http.Dir(*htmlPath)))
		http.HandleFunc("/shorten", ws.Shorten)
		log.Println("Registering web server on port:", port)
		log.Fatal(http.ListenAndServe(":"+*port, nil))
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
