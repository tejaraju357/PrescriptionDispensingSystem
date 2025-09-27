package cache

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var Rdb *redis.Client
var Ctx = context.Background()

func InitRedis() {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     "redis-17345.crce179.ap-south-1-1.ec2.redns.redis-cloud.com:17345",
		Password: "wRSHn2U5aLfY0zyn0wQYslBClXOUoWjp", // Replace with your Redis password
		DB:       0,
	})
	Pong, err := Rdb.Ping(Ctx).Result()
	if err != nil {
		fmt.Println("Error connecting to Redis:", err)
		return
	}
	fmt.Printf("Connected to Redis: %s\n", Pong)
}
