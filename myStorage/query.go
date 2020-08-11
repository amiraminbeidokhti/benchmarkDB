package myStorage

import (
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/amiraminbeidokhti/benchmarkDB/data"

	"github.com/amiraminbeidokhti/benchmarkDB/redis"
	redigo "github.com/gomodule/redigo/redis"
)

type MyStorage struct {
	// THIS IS LIKE TEST TABLE
	db map[int]string
}

var (
	host = os.Getenv("REDIS_HOST_MYSTORAGE")
	port = os.Getenv("REDIS_PORT_MYSTORAGE")

	numOfData, _    = strconv.Atoi(os.Getenv("NUM_OF_DATA"))
	lengthOfData, _ = strconv.Atoi(os.Getenv("LENGTH_OF_DATA"))

	replica redis.RedisDB
	mu      sync.Mutex
)

func (db *MyStorage) CreateConn() {
	db.db = make(map[int]string)
	replica.Pool = createReplicaPool()
}

func (db *MyStorage) Insert() {
	for i := 0; i < numOfData; i++ {
		prepare := data.RandString(lengthOfData)
		db.db[i] = prepare
		go insertReplica(&replica, i, prepare)
	}
}

func (db *MyStorage) Select() {
	temp := 0
	for k, _ := range db.db {
		temp += k
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
	_, err := conn.Do("HSET", "test", param[0], param[1])
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
