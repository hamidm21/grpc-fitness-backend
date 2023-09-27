package redis

import (
	"github.com/go-redis/redis"
	"gitlab.com/mefit/mefit-server/utils/assert"
	"gitlab.com/mefit/mefit-server/utils/initializer"
)

var Client *redis.Client

type manager struct{}

func (manager) Initialize() func() {
	Client = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})
	_, err := Client.Ping().Result()
	assert.Nil(err)

	return nil
}

func init() {
	initializer.Register(manager{})
}
