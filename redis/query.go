package redis

import (
	"fmt"
	"os"
	"strconv"

	"github.com/amiraminbeidokhti/benchmarkDB/data"

	"github.com/gomodule/redigo/redis"
)

type RedisDB struct {
	Pool *redis.Pool
}

var (
	host            = os.Getenv("REDIS_HOST")
	port            = os.Getenv("REDIS_PORT")
	sync            = os.Getenv("REDIS_SYNC")
	numOfReplica, _ = strconv.Atoi(os.Getenv("NUM_OF_REPLICA"))

	numOfData, _    = strconv.Atoi(os.Getenv("NUM_OF_DATA"))
	lengthOfData, _ = strconv.Atoi(os.Getenv("LENGTH_OF_DATA"))
)

func (db *RedisDB) CreateConn() error {

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
	return nil
}

func (db *RedisDB) Insert() {
	if sync == "1" {
		insertSync(db)
	} else {
		insertAsync(db)
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
	if sync == "1" {
		deleteSync(db)
	} else {
		deleteAsync(db)
	}
}

func insertAsync(db *RedisDB) {
	conn := db.Pool.Get()
	defer conn.Close()
	for i := 0; i < numOfData; i++ {
		s := data.RandString(lengthOfData)
		_, err := conn.Do("HSET", "test", i, s)
		if err != nil {
			fmt.Errorf(err.Error())
		}
	}
}

func insertSync(db *RedisDB) {
	conn := db.Pool.Get()
	defer conn.Close()
	for i := 0; i < numOfData; i++ {
		s := data.RandString(lengthOfData)
		_, err := conn.Do("HSET", "test", i, s)
		if err != nil {
			fmt.Errorf(err.Error())
		}
		_, err = conn.Do("wait", numOfReplica, 0)
		if err != nil {
			fmt.Errorf(err.Error())
		}
	}
}

func deleteAsync(db *RedisDB) {
	conn := db.Pool.Get()
	defer conn.Close()
	_, err := conn.Do("DEL", "test")
	if err != nil {
		fmt.Errorf(err.Error())
	}
}

func deleteSync(db *RedisDB) {
	conn := db.Pool.Get()
	defer conn.Close()
	_, err := conn.Do("DEL", "test")
	if err != nil {
		fmt.Errorf(err.Error())
	}
	_, err = conn.Do("wait", numOfReplica, 0)
	if err != nil {
		fmt.Errorf(err.Error())
	}
}
