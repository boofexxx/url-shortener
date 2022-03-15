package store

import (
	"context"
	"time"

	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
)

type Store struct {
	rdb   *redis.Client
	cache cache.Cache
	ctx   context.Context
}

func NewStore() *Store {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	return &Store{
		rdb: rdb,
		cache: *cache.New(&cache.Options{
			Redis:      rdb,
			LocalCache: cache.NewTinyLFU(5000, 5*time.Minute),
		}),
		ctx: context.Background(),
	}
}

// Ping checks if connection is established
func (s *Store) Ping() error {
	_, err := s.rdb.Ping(s.ctx).Result()
	return err
}

// Add maps key with value
func (s *Store) Add(key, value string) error {
	err := s.cache.Set(&cache.Item{
		Ctx:   s.ctx,
		Key:   key,
		Value: value,
		TTL:   5 * time.Minute,
	})
	return err
}

// Get returns value mapped with key
func (s *Store) Get(key string) (string, error) {
	value := ""
	err := s.cache.Get(s.ctx, key, &value)
	if err != nil {
		return "", err
	}
	return value, nil
}
