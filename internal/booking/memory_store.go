package booking

type MemoryStore struct {
	// seats --> booking
	bookings map[string]Booking
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		bookings: map[string]Booking{},
	}
}

func (s *MemoryStore) Book(b Booking) error {
	_, exists := s.bookings[b.SeatID]

	if exists {
		return ErrSeatAlreadyBooked
	}

	s.bookings[b.SeatID] = b
	return nil
}

func (s *MemoryStore) ListBookings(movieID string) []Booking {
	var bookings []Booking

	for _, booking := range s.bookings {
		if booking.MovieID == movieID {
			bookings = append(bookings, booking)
		}
	}

	return bookings
}
