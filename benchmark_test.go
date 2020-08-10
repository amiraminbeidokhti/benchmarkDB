package BenchmarkDB

import (
	"testing"

	"github.com/amiraminbeidokhti/benchmarkDB/mySQL"
	"github.com/amiraminbeidokhti/benchmarkDB/myStorage"
	"github.com/amiraminbeidokhti/benchmarkDB/redis"
)

// MYSQL BENCHMARL
var mySQLDb = mySQL.MySQLDb{}

func prepareMySQLDb() {
	mySQLDb.CreateMySQLCon()
	deleteDB(&mySQLDb)
	insertDB(&mySQLDb)
}

func BenchmarkInsertMySQL(b *testing.B) {
	prepareMySQLDb()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		insertDB(&mySQLDb)
	}
}

func BenchmarkSelectMySQL(b *testing.B) {
	prepareMySQLDb()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		selectDB(&mySQLDb)
	}
}

func BenchmarkDeleteMySQL(b *testing.B) {
	prepareMySQLDb()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		deleteDB(&mySQLDb)
	}
}

// // POSTGRESQL BENCHMARK
// var postgresDb = postgreSQL.PostgreSQLDb{}

// func preparePostgresDb() {
// 	postgresDb.CreatePostgreSQLCon()
// 	deleteDB(&postgresDb)
// 	insertDB(&postgresDb)
// }

// func BenchmarkInsertPostgres(b *testing.B) {
// 	preparePostgresDb()
// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		insertDB(&postgresDb)
// 	}
// }

// func BenchmarkSelectPostgres(b *testing.B) {
// 	preparePostgresDb()
// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		selectDB(&postgresDb)
// 	}
// }

// func BenchmarkDeletePostgres(b *testing.B) {
// 	preparePostgresDb()
// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		deleteDB(&postgresDb)
// 	}
// }

// REDIS
var redisDB = redis.RedisDB{}

func prepareRedisDb() {
	redisDB.CreateRedisPool()
	deleteDB(&redisDB)
	insertDB(&redisDB)
}

func BenchmarkInsertRedis(b *testing.B) {
	prepareRedisDb()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		insertDB(&redisDB)
	}
}

func BenchmarkSelectRedis(b *testing.B) {
	prepareRedisDb()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		selectDB(&redisDB)
	}
}

func BenchmarkDeleteRedis(b *testing.B) {
	prepareRedisDb()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		deleteDB(&redisDB)
	}
}

//	MYSTORAGE
var ms = myStorage.MyStorage{}

func prepareMyStorage() {
	ms.CreateDB()
	deleteDB(&ms)
	insertDB(&ms)
}

func BenchmarkInsertMyStorage(b *testing.B) {
	prepareMyStorage()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		insertDB(&ms)
	}
}

func BenchmarkSelectMystorage(b *testing.B) {
	prepareMyStorage()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		selectDB(&ms)
	}
}

func BenchmarkDeleteMyStorage(b *testing.B) {
	prepareMyStorage()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		deleteDB(&ms)
	}
}
