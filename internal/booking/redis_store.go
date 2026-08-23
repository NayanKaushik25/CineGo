package booking

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const defaultHoldTTL = 2 * time.Minute

// Implements our Redis solution
type RedisStore struct {
	rdb *redis.Client
}

func NewRedisStore(rdb *redis.Client) *RedisStore {
	return &RedisStore{rdb: rdb}
}

func sessionKey(id string) string {
	return fmt.Sprintf("session:%s", id)
}

func (s *RedisStore) Book(b Booking) (Booking, error) {
	session, err := s.hold(b)
	if err != nil {
		return Booking{}, err
	}
	log.Printf("Session booked %v", session)
	return session, nil
}

func (s *RedisStore) ListBookings(movieID string) []Booking {
	ctx := context.Background()
	pattern := fmt.Sprintf("seat:%s:*", movieID)
	var sessions []Booking
	iter := s.rdb.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		val, err := s.rdb.Get(ctx, iter.Val()).Result()
		if err != nil {
			continue
		}
		session, err := parseSession(val)
		if err != nil {
			continue
		}
		sessions = append(sessions, session)
	}
	return sessions
}

func parseSession(val string) (Booking, error) {
	var data Booking
	if err := json.Unmarshal([]byte(val), &data); err != nil {
		return Booking{}, err
	}
	return Booking{
		ID:        data.ID,
		MovieID:   data.MovieID,
		SeatID:    data.SeatID,
		UserID:    data.UserID,
		Status:    data.Status,
		ExpiresAt: data.ExpiresAt,
	}, nil
}

func (s *RedisStore) hold(b Booking) (Booking, error) {
	id := uuid.New().String()
	now := time.Now()
	ctx := context.Background()
	key := fmt.Sprintf("seat:%s:%s", b.MovieID, b.SeatID)
	b.ID = id
	val, err := json.Marshal(b)
	if err != nil {
		return Booking{}, err
	}
	res := s.rdb.SetArgs(ctx, key, val, redis.SetArgs{
		Mode: "NX",
		TTL:  defaultHoldTTL,
	})
	if err := res.Err(); err != nil {
		return Booking{}, err
	}
	ok := res.Val() == "OK"
	if !ok {
		return Booking{}, ErrSeatAlreadyTaken
	}
	if err := s.rdb.Set(ctx, sessionKey(id), key, defaultHoldTTL).Err(); err != nil {
		return Booking{}, err
	}
	return Booking{
		ID:        id,
		MovieID:   b.MovieID,
		SeatID:    b.SeatID,
		UserID:    b.UserID,
		Status:    "held",
		ExpiresAt: now.Add(defaultHoldTTL),
	}, nil
}
