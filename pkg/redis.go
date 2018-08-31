package root

import (
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
)

func InitRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis.address"),
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.db"),
	})

	_, err := client.Ping().Result()
	checkErr(err)

	return client
}