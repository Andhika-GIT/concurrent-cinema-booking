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

// seat:{movieID}:{seatID} -> session JSON (TTL = held, no TTL = confirmed)
// session:{sessionID} -> seat key ( reverse lookup )
type RedisStore struct {
	rdb *redis.Client
}

func NewRedisStore(rdb *redis.Client) *RedisStore {
	return &RedisStore{
		rdb: rdb,
	}
}

// sessionKey builds the reverse-lookup key for a session
func sessionKey(id string) string {
	return fmt.Sprintf("session:%s", id)
}

func (r *RedisStore) Book(b Booking) error {
	session, err := r.Hold(b)

	if err != nil {
		return err
	}

	log.Printf("session booked %v", session)

	return nil
}

func (r *RedisStore) ListBookings(movieID string) []Booking {
	return []Booking{}
}

func (r *RedisStore) Hold(b Booking) (Booking, error) {
	id := uuid.New().String()
	now := time.Now()
	key := fmt.Sprintf("seat:%s:%s", b.MovieID, b.SeatID)
	ctx := context.Background()

	b.ID = id
	val, _ := json.Marshal(b)

	res := r.rdb.SetArgs(ctx, key, val, redis.SetArgs{
		Mode: "NX",
		TTL:  defaultHoldTTL,
	})

	ok := res.Val() == "OK"

	if !ok {
		return Booking{}, ErrSeatAlreadyBooked
	}

	r.rdb.Set(ctx, sessionKey(id), key, defaultHoldTTL)

	return Booking{
		ID:        id,
		MovieID:   b.MovieID,
		SeatID:    b.SeatID,
		UserID:    b.UserID,
		Status:    "held",
		ExpiresAt: now.Add(defaultHoldTTL),
	}, nil
}
