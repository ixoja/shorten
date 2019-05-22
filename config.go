package main

import (
	"github.com/alecthomas/kingpin"
	"github.com/ixoja/shorten/internal/restapi"
	"log"
)

type Config struct {
	mode       string
	htmlPath   string
	port       string
	apiURL     string
	apiAddress string
	apiConfig  restapi.Config
}

func (c *Config) WithFlags() *Config {
	kingpin.Flag("mode", "string value: webserver or apiserver").StringVar(&c.mode)
	kingpin.Flag("html", "path to html dir").StringVar(&c.htmlPath)
	kingpin.Flag("port", "webserver http port").StringVar(&c.port)
	kingpin.Flag("api_url", "web api url").StringVar(&c.apiURL)
	kingpin.Flag("api_address", "web api serve address").StringVar(&c.apiConfig.Address)
	c.apiConfig.InsecureHTTP = true
	err := c.apiConfig.Parse()
	if err != nil {
		log.Fatal("Failed to parse config", err)
	}
	//c.apiConfig.Address = c.apiAddress

	return c
}
