package infra

import (
	"context"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

//IMemCash interface to store key values
type IMemCash interface {
	GetKey(key string) (string, error)
	SetKey(key string, value interface{}, expiration time.Duration) error
	DelKey(key string) (int64, error)
}

type redisMem struct {
	cl *redis.Client
}

func (r *redisMem) GetKey(key string) (string, error) {
	return r.cl.Get(context.Background(), key).Result()
}

func (r *redisMem) SetKey(key string, value interface{}, expiration time.Duration) error {
	return r.cl.Set(context.Background(), key, value, expiration).Err()
}

func (r *redisMem) DelKey(key string) (int64, error) {
	return r.cl.Del(context.Background(), key).Result()
}

//NewMemCash init redis connection
func NewMemCash() (IMemCash, error) {
	//Initializing redis
	dsn := os.Getenv("REDIS_DSN")
	if len(dsn) == 0 {
		dsn = "localhost:6379"
	}
	client := redis.NewClient(&redis.Options{
		Addr: dsn,
	})
	//Ping Redis
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}
	return &redisMem{cl: client}, nil
}
