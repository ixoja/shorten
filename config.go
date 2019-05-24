package main

import (
	"github.com/alecthomas/kingpin"
)

type Config struct {
	mode     string
	htmlPath string
	port     string
	webURL   string
	apiURL   string
}

func (c *Config) WithFlags() *Config {
	kingpin.Flag("mode", "string value: webserver or apiserver").StringVar(&c.mode)
	kingpin.Flag("html", "path to html dir").StringVar(&c.htmlPath)
	kingpin.Flag("port", "webserver http port").StringVar(&c.port)
	kingpin.Flag("web_url", "web server url").StringVar(&c.webURL)
	kingpin.Flag("api_url", "backend api url").StringVar(&c.apiURL)
	kingpin.Parse()

	return c
}
