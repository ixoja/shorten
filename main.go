package main

import (
	"log"

	"github.com/ixoja/shorten/internal/handler"
	"github.com/ixoja/shorten/internal/restapi"
)

func main() {
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
