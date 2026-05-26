package cache

import (
	"fmt"

	"github.com/boilerplate/internal/config"
	"github.com/redis/go-redis/v9"
)

func New(cfg config.RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
	})
	return client, nil
}
