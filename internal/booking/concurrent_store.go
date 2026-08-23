package booking

import "sync"

type ConcurrentStore struct {
	// seat -> booking mapping
	bookings map[string]Booking
	sync.RWMutex
}

func NewConcurrentStore() *ConcurrentStore {
	return &ConcurrentStore{
		bookings: map[string]Booking{},
	}
}

func (s *ConcurrentStore) Book(b Booking) (Booking, error) {
	s.Lock()
	defer s.Unlock()
	if _, exists := s.bookings[b.SeatID]; exists {
		return Booking{}, ErrSeatAlreadyTaken
	}
	b.Status = "held"
	s.bookings[b.SeatID] = b
	return b, nil
}

func (s *ConcurrentStore) ListBookings(movieID string) []Booking {
	s.RLock()
	defer s.RUnlock()
	var result []Booking
	for _, b := range s.bookings {
		if b.MovieID == movieID {
			result = append(result, b)
		}
	}
	return result
}
