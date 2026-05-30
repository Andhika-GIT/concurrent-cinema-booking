package booking

import "sync"

type ConcurrentStore struct {
	// seats --> booking
	bookings map[string]Booking
	sync.RWMutex
}

func NewConcurrentStore() *ConcurrentStore {
	return &ConcurrentStore{
		bookings: map[string]Booking{},
	}
}

func (s *ConcurrentStore) Book(b Booking) error {
	s.Lock()
	defer s.Unlock()

	_, exists := s.bookings[b.SeatID]

	if exists {
		return ErrSeatAlreadyBooked
	}

	s.bookings[b.SeatID] = b
	return nil
}

func (s *ConcurrentStore) ListBookings(movieID string) []Booking {
	s.RLock()
	defer s.RUnlock()

	var bookings []Booking

	for _, booking := range s.bookings {
		if booking.MovieID == movieID {
			bookings = append(bookings, booking)
		}
	}

	return bookings
}
