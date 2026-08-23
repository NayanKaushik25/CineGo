package booking

type MemoryStore struct {
	// seat -> booking mapping
	bookings map[string]Booking
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		bookings: map[string]Booking{},
	}
}

func (s *MemoryStore) Book(b Booking) (Booking, error) {
	if _, exists := s.bookings[b.SeatID]; exists {
		return Booking{}, ErrSeatAlreadyTaken
	}
	b.Status = "held"
	s.bookings[b.SeatID] = b
	return b, nil
}

func (s *MemoryStore) ListBookings(movieID string) []Booking {
	var result []Booking
	for _, b := range s.bookings {
		if b.MovieID == movieID {
			result = append(result, b)
		}
	}
	return result
}
