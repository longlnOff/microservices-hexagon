package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type Redis struct {
	client *redis.Client
}

func NewValkeyClient(address string, password string, db int) *Redis {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       db,
	})

	return &Redis{
		client: redisClient,
	}
}

// Set stores the value in the redis database
func (r *Redis) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl).Err()
}

// Get retrieves the value from the redis database
func (r *Redis) Get(ctx context.Context, key string) ([]byte, error) {
	res, err := r.client.Get(ctx, key).Result()
	bytes := []byte(res)
	return bytes, err
}

// Delete removes the value from the redis database
func (r *Redis) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

// DeleteByPrefix removes the value from the redis database with the given prefix
func (r *Redis) DeleteByPrefix(ctx context.Context, prefix string) error {
	var cursor uint64
	var keys []string

	for {
		var err error
		keys, cursor, err = r.client.Scan(ctx, cursor, prefix, 100).Result()
		if err != nil {
			return err
		}

		for _, key := range keys {
			err := r.client.Del(ctx, key).Err()
			if err != nil {
				return err
			}
		}

		if cursor == 0 {
			break
		}
	}

	return nil
}

// Close closes the connection to the redis database
func (r *Redis) Close() error {
	return r.client.Close()
}
