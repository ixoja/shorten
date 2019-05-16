package main

import (
	"github.com/ixoja/shorten/internal/handler"
	"log"

	"github.com/ixoja/shorten/internal/restapi"
)

var apiConfig restapi.Config

func main() {
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
