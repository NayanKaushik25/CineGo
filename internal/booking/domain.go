package booking

import (
	"errors"
	"time"
)

var (
	ErrSeatAlreadyTaken = errors.New("Seat is already taken")
	ErrNotImplemented   = errors.New("feature not implemented")
)

type Booking struct {
	ID        string
	MovieID   string
	SeatID    string
	UserID    string
	Status    string
	ExpiresAt time.Time
}

type BookingStore interface {
	Book(b Booking) (Booking, error)
	ListBookings(MovieID string) []Booking
}
