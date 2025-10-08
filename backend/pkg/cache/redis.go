package cache

import (
	"context"
	"log"

	"magicchat/pkg/config"
	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func ConnectRedis(cfg *config.Config) (*redis.Client, error) {
	opt, err := redis.ParseURL(cfg.Redis.URL)
	if err != nil {
		return nil, err
	}

	if cfg.Redis.Password != "" {
		opt.Password = cfg.Redis.Password
	}

	RedisClient = redis.NewClient(opt)

	// Test connection
	ctx := context.Background()
	_, err = RedisClient.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	log.Println("Connected to Redis successfully")
	return RedisClient, nil
}

func DisconnectRedis() error {
	if RedisClient != nil {
		return RedisClient.Close()
	}
	return nil
}
