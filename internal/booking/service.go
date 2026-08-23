package booking

import "context"

type Service struct {
	store BookingStore
}

func NewService(store BookingStore) *Service {
	return &Service{store}
}

func (s *Service) Book(b Booking) (Booking, error) {
	return s.store.Book(b)
}

func (s *Service) ListBookings(movieID string) []Booking {
	return s.store.ListBookings(movieID)
}

func (s *Service) ListBooking(movieID string) []Booking {
	return s.ListBookings(movieID)
}

func (s *Service) ConfirmSeat(_ context.Context, _ string, _ string) (Booking, error) {
	return Booking{}, ErrNotImplemented
}

func (s *Service) ReleaseSeat(_ context.Context, _ string, _ string) error {
	return ErrNotImplemented
}
