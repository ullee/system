package redis

import (
	"github.com/go-redis/redis"
	"time"
)

const (
	REDIS_ADDR = "localhost"
	REDIS_PORT = "6379"
	REDIS_PASS = ""
	REDIS_DB   = 0
)

type Redis struct {
	client *redis.Client
}

func (r *Redis) Init() {
	r.setClient()
}

func (r *Redis) setClient() {
	r.client = redis.NewClient(&redis.Options{
		Addr:     REDIS_ADDR+":"+REDIS_PORT,
		Password: REDIS_PASS,
		DB:       REDIS_DB,
	})
}

func (r *Redis) Ping() (string, error) {
	pong, err := r.client.Ping().Result()
	return pong, err
}

func (r *Redis) Set(key string, value interface{}, expiration time.Duration) error {
	err := r.client.Set(key, value, expiration).Err()
	return err
}

func (r *Redis) Get(key string) (string, error) {
	val, err := r.client.Get(key).Result()
	return val, err
}