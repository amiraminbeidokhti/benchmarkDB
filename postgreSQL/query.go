package postgreSQL

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	_ "github.com/lib/pq"
)

type PostgreSQLDb struct {
	db *sql.DB
}

var (
	host  = os.Getenv("POSTGRES_HOST")
	port  = os.Getenv("POSTGRES_PORT")
	user  = os.Getenv("POSTGRES_USER")
	pass  = os.Getenv("POSTGRES_PASSWORD")
	dbase = os.Getenv("POSTGRES_DBNAME")
)

func (db *PostgreSQLDb) CreatePostgreSQLCon() {
	port, err := strconv.Atoi(port)
	if err != nil {
		fmt.Errorf(err.Error())
	}
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, pass, dbase)
	ps, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("Failed to open a DB connection: ", err)
	}
	db.db = ps

	err = db.db.Ping()
	if err != nil {
		fmt.Errorf("PostgreSQL db is not connected")
		fmt.Println(err.Error())
	}
	createTable(db)
}

func (db *PostgreSQLDb) Insert() {
	for i := 0; i < 1000; i++ {
		stmtIns, err := db.db.Prepare("INSERT INTO test VALUES (DEFAULT, $1, $2, $3);")
		defer stmtIns.Close()
		if err != nil {
			panic(err)
		}
		_, err = stmtIns.Exec(i, i, i)
		if err != nil {
			panic(err)
		}
	}
}

func (db *PostgreSQLDb) Select() {
	results, err := db.db.Query("SELECT * FROM test;")
	if err != nil {
		panic(err.Error())
	}
	results.Close()
}

func (db *PostgreSQLDb) Delete() {
	del, err := db.db.Query("DELETE FROM test;")
	if err != nil {
		fmt.Errorf(err.Error())
	}
	del.Close()
}

func createTable(db *PostgreSQLDb) {
	result, err := db.db.Query("CREATE TABLE IF NOT EXISTS test (id serial primary key, f1 int, f2 int, f3 int);")
	if err != nil {
		panic(err.Error())
	}
	result.Close()
}
