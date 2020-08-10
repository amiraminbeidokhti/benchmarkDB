package mySQL

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type MySQLDb struct {
	db *sql.DB
}

var (
	user  = os.Getenv("MYSQL_USER")
	pass  = os.Getenv("MYSQL_PASSWORD")
	host  = os.Getenv("MYSQL_HOST")
	port  = os.Getenv("MYSQL_PORT")
	dbase = os.Getenv("MYSQL_DBNAME")
)

func (db *MySQLDb) CreateMySQLCon() {
	s := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, pass, host, port, dbase)
	ms, err := sql.Open("mysql", s)
	if err != nil {
		fmt.Errorf(err.Error())
	}
	db.db = ms
	// make sure connection is available
	err = db.db.Ping()
	if err != nil {
		fmt.Errorf("MySQL db is not connected")
		fmt.Println(err.Error())
	}
	createTable(db)
}

func (db *MySQLDb) Insert() {
	for i := 0; i < 1000; i++ {
		stmtIns, err := db.db.Prepare("INSERT INTO test VALUES (null, ?, ?, ?);")
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

func (db *MySQLDb) Select() {
	results, err := db.db.Query("SELECT * FROM test;")
	if err != nil {
		panic(err.Error())
	}
	results.Close()
}

func (db *MySQLDb) Delete() {
	del, err := db.db.Query("DELETE FROM test;")
	if err != nil {
		fmt.Errorf(err.Error())
	}
	del.Close()
}

func createTable(db *MySQLDb) {
	result, err := db.db.Query("CREATE TABLE IF NOT EXISTS test (id int(6) primary key auto_increment, f1 int(6), f2 int(6), f3 int(6));")
	if err != nil {
		panic(err.Error())
	}
	result.Close()
}
