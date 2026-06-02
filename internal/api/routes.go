package api

import "net/http"

func SetupRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /movies", listMovies)
}
