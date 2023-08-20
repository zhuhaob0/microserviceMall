package setup

import (
	conf "final-design/pkg/config"
	"fmt"
	"log"

	"github.com/go-redis/redis"
)

// 初始化redis
func InitRedis() {
	client := redis.NewClient(&redis.Options{
		Addr:     conf.Redis.Host,
		Password: conf.Redis.Password,
		DB:       conf.Redis.Db,
	})

	_, err := client.Ping().Result()
	if err != nil {
		log.Printf("Connect redis failed, Error: %v", err)
	}
	conf.Redis.RedisConn = client
	fmt.Println("Redis 连接成功")
}
