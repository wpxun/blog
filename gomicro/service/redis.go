package service

import (
	"github.com/go-redis/redis"
)

func GetRedis() *redis.Client {
	return redis.NewClient(&redis.Options {
		Addr:         Conf.Redis.Address,
		Password:     Conf.Redis.Password,
		DB:           Conf.Redis.Database,
		Network:      Conf.Redis.Network,
		DialTimeout:  Conf.Redis.DialTimeout,
	})
}