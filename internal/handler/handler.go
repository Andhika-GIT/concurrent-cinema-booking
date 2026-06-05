package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Andhika-GIT/concurrent-cinema-booking/internal/booking"
)

type BookingHandlerInterface interface {
	ListMovies(w http.ResponseWriter, r *http.Request)
	ListSeats(w http.ResponseWriter, r *http.Request)
	HoldSeat(w http.ResponseWriter, r *http.Request)
	ConfirmSession(w http.ResponseWriter, r *http.Request)
	ReleaseSession(w http.ResponseWriter, r *http.Request)
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

func (b *BookingHandler) ListSeats(w http.ResponseWriter, r *http.Request) {
	movieID := r.PathValue("movieID")

	bookings := b.service.ListBookings(movieID)

	seats := make([]seatInfo, 0, len(bookings))
	for _, b := range bookings {
		seats = append(seats, seatInfo{
			SeatID:    b.SeatID,
			UserID:    b.UserID,
			Booked:    true,
			Confirmed: b.Status == "confirmed",
		})
	}

	writeJSON(w, http.StatusOK, seats)
}

func (b *BookingHandler) HoldSeat(w http.ResponseWriter, r *http.Request) {
	movieID := r.PathValue("movieID")
	seatID := r.PathValue("seatID")

	var req holdRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		log.Println(err)
		writeJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	session, err := b.service.Book(booking.Booking{
		UserID:  req.UserID,
		SeatID:  seatID,
		MovieID: movieID,
	})

	if err != nil {
		writeJSON(w, http.StatusConflict, err.Error())
		log.Println(err)
		return
	}

	writeJSON(w, http.StatusCreated, holdResponse{
		SeatID:    seatID,
		MovieID:   movieID,
		SessionID: session.ID,
		ExpiresAt: session.ExpiresAt.Format(time.RFC3339),
	})
}

func (b *BookingHandler) ConfirmSession(w http.ResponseWriter, r *http.Request) {
	sessionID := r.PathValue("sessionID")

	err := b.service.ConfirmSession(sessionID)

	if err != nil {
		log.Println(err)
		writeJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusNoContent, "")
}

func (b *BookingHandler) ReleaseSession(w http.ResponseWriter, r *http.Request) {
	sessionID := r.PathValue("sessionID")

	err := b.service.ReleaseSession(sessionID)

	if err != nil {
		log.Println(err)
		writeJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusNoContent, "")
}

// ================== HELPER ================== //

type movieResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Rows        int    `json:"rows"`
	SeatsPerRow int    `json:"seats_per_row"`
}

type holdResponse struct {
	SessionID string `json:"session_id"`
	MovieID   string `json:"movieID"`
	SeatID    string `json:"seat_id"`
	ExpiresAt string `json:"expires_at"`
}

type seatInfo struct {
	SeatID    string `json:"seat_id"`
	UserID    string `json:"user_id"`
	Booked    bool   `json:"booked"`
	Confirmed bool   `json:"confirmed"`
}

type holdRequest struct {
	UserID string `json:"user_id"`
}

var movies = []movieResponse{
	{ID: "inception", Title: "Inception", Rows: 5, SeatsPerRow: 8},
	{ID: "dune", Title: "Dune: Part Two", Rows: 4, SeatsPerRow: 6},
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)

}
