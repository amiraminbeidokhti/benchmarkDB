package BenchmarkDB

import (
	"testing"

	"github.com/amiraminbeidokhti/benchmarkDB/mySQL"
	"github.com/amiraminbeidokhti/benchmarkDB/myStorage"
	"github.com/amiraminbeidokhti/benchmarkDB/postgreSQL"
	"github.com/amiraminbeidokhti/benchmarkDB/redis"
)

type memory interface {
	CreateConn()
	Delete()
	Insert()
}

func prepareMemory(m memory) {
	m.CreateConn()
	m.Delete()
	m.Insert()
}

func benchmarkInsert(b *testing.B, s storage) {
	prepareMemory(s.(memory))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		insertDB(s)
	}
}

func benchmarkSelect(b *testing.B, s storage) {
	prepareMemory(s.(memory))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		selectDB(s)
	}
}

func benchmarkDelete(b *testing.B, s storage) {
	prepareMemory(s.(memory))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		deleteDB(s)
	}
}

// MYSQL BENCHMARK
var mySQLDb = &mySQL.MySQLDb{}

func BenchmarkInsertMySQL(b *testing.B) {
	benchmarkInsert(b, mySQLDb)
}

func BenchmarkSelectMySQL(b *testing.B) {
	benchmarkSelect(b, mySQLDb)
}
func BenchmarkDeleteMySQL(b *testing.B) {
	benchmarkDelete(b, mySQLDb)
}

// POSTGRESQL BENCHMARK
var postgresDb = &postgreSQL.PostgreSQLDb{}

func BenchmarkInsertPostgres(b *testing.B) {
	benchmarkInsert(b, postgresDb)
}
func BenchmarkSelectPostgres(b *testing.B) {
	benchmarkSelect(b, postgresDb)
}
func BenchmarkDeletePostgres(b *testing.B) {
	benchmarkDelete(b, postgresDb)
}

// REDIS
var redisDB = &redis.RedisDB{}

func BenchmarkInsertRedis(b *testing.B) {
	benchmarkInsert(b, redisDB)
}
func BenchmarkSelectRedis(b *testing.B) {
	benchmarkSelect(b, redisDB)
}
func BenchmarkDeleteRedis(b *testing.B) {
	benchmarkDelete(b, redisDB)
}

//	MYSTORAGE
var ms = &myStorage.MyStorage{}

func BenchmarkInsertMyStorage(b *testing.B) {
	benchmarkInsert(b, ms)
}
func BenchmarkSelectMystorage(b *testing.B) {
	benchmarkSelect(b, ms)
}
func BenchmarkDeleteMyStorage(b *testing.B) {
	benchmarkDelete(b, ms)
}
