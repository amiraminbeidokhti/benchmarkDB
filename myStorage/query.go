package myStorage

import (
	"fmt"
	"os"
	"sync"

	"github.com/amiraminbeidokhti/benchmarkDB/redis"
	redigo "github.com/gomodule/redigo/redis"
)

type MyStorage struct {
	db map[int][]int
}

var (
	host = os.Getenv("REDIS_HOST_MYSTORAGE")
	port = os.Getenv("REDIS_PORT_MYSTORAGE")

	replica redis.RedisDB
	mu      sync.Mutex
)

func (db *MyStorage) CreateDB() {
	db.db = make(map[int][]int)
	replica.Pool = createReplicaPool()
}

func (db *MyStorage) Insert() {
	for i := 0; i < 1000; i++ {
		prepare := []int{i, i, i}
		db.db[i] = prepare
		go insertReplica(&replica, i, i, i, i)
	}
}

func (db *MyStorage) Select() {
	temp := 0
	for range db.db {
		temp++
	}
}

func (db *MyStorage) Delete() {
	for k, _ := range db.db {
		delete(db.db, k)
	}
	deleteReplica()
}

func createReplicaPool() *redigo.Pool {
	return &redigo.Pool{
		MaxIdle:   80,
		MaxActive: 12000,
		Dial: func() (redigo.Conn, error) {
			s := fmt.Sprintf("%s:%s", host, port)
			c, err := redigo.Dial("tcp", s)
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}
}

func insertReplica(db *redis.RedisDB, param ...interface{}) {
	mu.Lock()
	conn := db.Pool.Get()
	defer conn.Close()
	_, err := conn.Do("HSET", "test", "id", param[0], "f1", param[1], "f2", param[2], "f3", param[3])
	if err != nil {
		fmt.Errorf(err.Error())
	}
	mu.Unlock()
}

func deleteReplica() {
	mu.Lock()
	replica.Delete()
	mu.Unlock()
}
