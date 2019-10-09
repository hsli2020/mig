package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
)

func Open() *sql.DB {
	dbhost := getenv("DBHOST", "localhost")
	dbport := getenv("DBPORT", "")
	dbuser := getenv("DBUSER", "root")
	dbpass := getenv("DBPASS", "")
	dbname := getenv("DBNAME", "bte")

	if dbport == "" {
		dbport = "3306"
	}

	var dsn string
	if dbpass == "" {
		dsn = fmt.Sprintf("%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true", dbuser, dbhost, dbport, dbname)
	} else {
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true", dbuser, dbpass, dbhost, dbport, dbname)
	}

	db, err := sql.Open("mysql", dsn)
	check(err)

	err = db.Ping()
	check(err)

	return db
}

func getenv(name, defaultval string) string {
	val := os.Getenv(name)
	if val == "" {
		return defaultval
	}
	return val
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
