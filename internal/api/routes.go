package api

import (
	"net/http"

	"github.com/Andhika-GIT/concurrent-cinema-booking/internal/handler"
)

func SetupRoutes(mux *http.ServeMux, bookingHandler handler.BookingHandlerInterface) {
	mux.Handle("GET /", http.FileServer(http.Dir("static")))
	mux.HandleFunc("GET /movies", bookingHandler.ListMovies)
}
