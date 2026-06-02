package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Andhika-GIT/concurrent-cinema-booking/internal/booking"
)

type movieResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Rows        int    `json:"rows"`
	SeatsPerRow int    `json:"seats_per_row"`
}

type BookingHandlerInterface interface {
	ListMovies(w http.ResponseWriter, r *http.Request)
}

type BookingHandler struct {
	service booking.BookingStore
}

func NewBookingHandler(service booking.BookingStore) BookingHandlerInterface {
	return &BookingHandler{
		service: service,
	}
}

func (b *BookingHandler) ListMovies(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, movies)
}

// ================== HELPER FUNCTION ================== //

var movies = []movieResponse{
	{ID: "inception", Title: "Inception", Rows: 5, SeatsPerRow: 8},
	{ID: "dune", Title: "Dune: Part Two", Rows: 4, SeatsPerRow: 6},
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)

}
