package webserver

import (
	"log"
	"net/http"
)

type Server struct {
	HtmlPath string
	Port     string
}

func New(htmlPath, port string) *Server {
	return &Server{HtmlPath: htmlPath, Port: port}
}

func (w *Server) Start() {
	http.Handle("/", http.FileServer(http.Dir(w.HtmlPath)))
	http.HandleFunc("/shorten", shorten)
	log.Println("Registering web server on port:", w.Port)
	log.Fatal(http.ListenAndServe(":"+w.Port, nil))
}

func shorten(w http.ResponseWriter, r *http.Request) {

}
