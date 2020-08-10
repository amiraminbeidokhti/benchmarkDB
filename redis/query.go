package redis

import (
	"fmt"
	"os"

	"github.com/gomodule/redigo/redis"
)

type RedisDB struct {
	Pool *redis.Pool
}

var (
	host = os.Getenv("REDIS_HOST")
	port = os.Getenv("REDIS_PORT")
)

func (db *RedisDB) CreateRedisPool() {

	temp := &redis.Pool{
		MaxIdle:   80,
		MaxActive: 12000,
		Dial: func() (redis.Conn, error) {
			s := fmt.Sprintf("%s:%s", host, port)
			c, err := redis.Dial("tcp", s)
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}
	db.Pool = temp
}

func (db *RedisDB) Insert() {
	conn := db.Pool.Get()
	defer conn.Close()
	for i := 0; i < 1000; i++ {
		_, err := conn.Do("HSET", "test", "id", i, "f1", i, "f2", i, "f3", i)
		if err != nil {
			fmt.Errorf(err.Error())
		}
	}
}

func (db *RedisDB) Select() {
	conn := db.Pool.Get()
	defer conn.Close()
	_, err := conn.Do("HVALS", "test")
	if err != nil {
		fmt.Errorf(err.Error())
	}
}

func (db *RedisDB) Delete() {
	conn := db.Pool.Get()
	defer conn.Close()
	_, err := conn.Do("DEL", "test")
	if err != nil {
		fmt.Errorf(err.Error())
	}
}
