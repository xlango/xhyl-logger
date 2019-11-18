package conf

import (
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

var Client *redis.Client

type RedisClient struct {
}

func initRedis(db []int) {
	dbCode := 0
	if len(db) > 0 {
		dbCode = db[0]
	}

	Client = redis.NewClient(&redis.Options{
		Addr:         GlobalConfig.RedisHost,
		PoolSize:     GlobalConfig.RedisPoolSize,
		ReadTimeout:  time.Millisecond * time.Duration(GlobalConfig.RedisReadTimeout),
		WriteTimeout: time.Millisecond * time.Duration(GlobalConfig.RedisWriteTimeout),
		IdleTimeout:  time.Second * time.Duration(GlobalConfig.RedisIdleTimeout),
		DB:           dbCode,
	})

	_, err := Client.Ping().Result()
	if err != nil {
		panic("init redis error")
	} else {
		fmt.Println("init redis ok")
	}
}

func (this *RedisClient) Get(key string, db ...int) (string, bool) {
	initRedis(db)
	r, err := Client.Get(key).Result()
	if err != nil {
		return "", false
	}
	return r, true
}

func (this *RedisClient) SetExpTime(key string, val interface{}, expTime int32, db ...int) {
	initRedis(db)
	Client.Set(key, val, time.Duration(expTime)*time.Second)
}

func (this *RedisClient) Set(key string, val interface{}) {

}
