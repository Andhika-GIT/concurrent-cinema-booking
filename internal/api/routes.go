package api

import "net/http"

func SetupRoutes(mux *http.ServeMux) {
	mux.Handle("GET /", http.FileServer(http.Dir("static")))
	mux.HandleFunc("GET /movies", listMovies)
}
