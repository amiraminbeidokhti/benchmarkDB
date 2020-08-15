package mySQL

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	"github.com/amiraminbeidokhti/benchmarkDB/data"
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

	numOfData, _    = strconv.Atoi(os.Getenv("NUM_OF_DATA"))
	lengthOfData, _ = strconv.Atoi(os.Getenv("LENGTH_OF_DATA"))
)

func (db *MySQLDb) CreateConn() error {
	s := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, pass, host, port, dbase)
	ms, err := sql.Open("mysql", s)
	if err != nil {
		return err
	}
	db.db = ms
	// make sure connection is available
	err = db.db.Ping()
	if err != nil {
		return err
	}
	if host == "mySqlServer" {
		createTable(db)
	} else {
		createTableSync(db)
	}
	return nil
}

func (db *MySQLDb) Insert() {
	for i := 0; i < numOfData; i++ {
		s := data.RandString(lengthOfData)
		results, err := db.db.Query(fmt.Sprintf(`INSERT INTO test VALUES (null, "%v");`, s))
		if err != nil {
			panic(err.Error())
		}
		results.Close()
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
		panic(err.Error())
	}
	del.Close()
}

func createTableSync(db *MySQLDb) {
	q := fmt.Sprintf("CREATE TABLE IF NOT EXISTS test (id int(10) primary key auto_increment, f1 varchar(%v)) ENGINE=NDBCLUSTER;", lengthOfData)
	result, err := db.db.Query(q)
	if err != nil {
		panic(err.Error())
	}
	result.Close()
}

func createTable(db *MySQLDb) {
	q := fmt.Sprintf("CREATE TABLE IF NOT EXISTS test (id int(10) primary key auto_increment, f1 varchar(%v));", lengthOfData)
	result, err := db.db.Query(q)
	if err != nil {
		panic(err.Error())
	}
	result.Close()
}
