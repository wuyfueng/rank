package redis_wrapper

import (
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

var (
	rdb *redis.Client
)

func Init() {
	host := viper.GetString("redis.host")
	port := viper.GetInt("redis.port")
	password := viper.GetString("redis.password")
	log.Println("redis.host", host)
	log.Println("redis.port", port)
	log.Println("redis.password", password)

	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: password,
	})
}

// RegisterRdb 测试用
func RegisterRdb(host string, port int, password string) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: password,
	})
}

func Rdb() *redis.Client {
	return rdb
}
