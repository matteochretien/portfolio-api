package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"os"
)

type Client struct {
	Client *redis.Client
}

func NewClient() (*Client, error) {
	redisUri, err := redis.ParseURL(os.Getenv("REDIS_URI"))
	if err != nil {
		return nil, err

	}
	redisClient := redis.NewClient(redisUri)
	_, err = redisClient.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return &Client{
		Client: redisClient,
	}, nil
}

func (r *Client) Close() error {
	return r.Client.Close()
}
