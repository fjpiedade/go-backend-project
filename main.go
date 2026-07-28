package main

import (
	"log"
	"net/http"
)

type server struct {
	addr string
}

// ServeHTTP implements http.Handler.
func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		switch r.URL.Path {
		case "/":
			w.Write([]byte("index page"))
			return
		case "/users":
			w.Write([]byte("user page"))
			return
		default:
			w.Write([]byte("404 page not found"))
			return
		}
	}
}

func main() {
	s := &server{
		addr: ":9090",
	}

	if err := http.ListenAndServe(s.addr, s); err != nil {
		log.Fatal(err)
	}
}
