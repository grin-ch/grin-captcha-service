package model

import (
	"github.com/go-redis/redis"
)

type Provider struct {
	cache *redis.Client
}

func RegistryProvider(addr, pass string, db int) (*Provider, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass,
		DB:       db,
	})
	_, err := rdb.Ping().Result()
	if err != nil {
		return nil, err
	}
	return &Provider{
		cache: rdb,
	}, nil
}
