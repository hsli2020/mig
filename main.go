package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

var db *sql.DB
var dbhost string
var dbport string
var dbuser string
var dbpass string
var dbname string
var migdir string
var migtab string

func main() {
	dbhost = getenv("DBHOST", "localhost")
	dbport = getenv("DBPORT", "")
	dbuser = getenv("DBUSER", "root")
	dbpass = getenv("DBPASS", "")
	dbname = getenv("DBNAME", "bte")
	migdir = getenv("MIGDIR", "migrations")
	migtab = getenv("MIGTAB", "migration_log")

	db = openDatabase(dbhost, dbport, dbuser, dbpass, dbname)
	defer db.Close()

	migrationInit()

	if len(os.Args) == 1 {
		fmt.Println("usage: mig new|up|down")
		return
	}

	if os.Args[1] == "new" {
		migrationNew()
	}

	if os.Args[1] == "up" {
		migrationUp()
	}

	if os.Args[1] == "down" {
		migrationDown()
	}
}

func migrationInit() {
	sql := "CREATE TABLE IF NOT EXISTS " + migtab + " (id varchar(20) not null primary key, action varchar(80) default '', run_at timestamp not null default CURRENT_TIMESTAMP) engine=InnoDB"
	db.Exec(sql)
}

func migrationNew() {
	if len(os.Args) < 3 {
		fmt.Println("usage: mig new <NAME>")
		return
	}

	action := os.Args[2]
	now := time.Now().Format("20060102-150405")

	// create migration up script
	filename := filepath.Join(migdir, now+"_"+action+".up.sql")
	file, err := os.Create(filename)
	check(err)
	file.WriteString("USE " + dbname + "\n")
	file.Close()
	fmt.Println("Created:", filename)

	// create migration down script
	filename = filepath.Join(migdir, now+"_"+action+".down.sql")
	file, err = os.Create(filename)
	check(err)
	file.WriteString("USE " + dbname + "\n")
	file.Close()
	fmt.Println("Created:", filename)
}

func migrationUp() {
	pattern := filepath.Join(migdir, "*.up.sql")

	files, err := filepath.Glob(pattern)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, file := range files {
		basename := filepath.Base(file)
		id := basename[0:15]
		action := basename[16:]

		sql := fmt.Sprintf("INSERT INTO %s set id='%s', action='%s'", migtab, id, action)
		_, err = db.Exec(sql)
		if err != nil {
			continue // next migration script
		}

		script, err := ioutil.ReadFile(file)
		if err != nil {
			continue // next migration script
		}

		_, err = db.Exec(string(script))
		if err != nil {
			continue // next migration script
		}

		fmt.Println("Running:", file)
	}
}

func migrationDown() {
	var id, action, sql string

	// fetch last migration
	sql = fmt.Sprintf("SELECT id, action FROM %s ORDER BY id DESC LIMIT 1", migtab)
	err := db.QueryRow(sql).Scan(&id, &action)
	check(err)

	// find the migration down script
	filename := filepath.Join(migdir, id+"_"+action[0:len(action)-7]+".down.sql")
	script, err := ioutil.ReadFile(filename)
	check(err)

	// run the migration script
	_, err = db.Exec(string(script))
	check(err)

	sql = fmt.Sprintf("DELETE FROM %s WHERE id='%s'", migtab, id)
	_, err = db.Exec(sql)
	check(err)

	fmt.Println("Running:", filename)
}

func getenv(name, defaultval string) string {
	val := os.Getenv(name)
	if val == "" {
		return defaultval
	}
	return val
}

func openDatabase(host, port, user, pass, dbname string) *sql.DB {
	if port == "" {
		port = "3306"
	}

	var dsn string
	if pass == "" {
		dsn = fmt.Sprintf("%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true", user, host, port, dbname)
	} else {
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true", user, pass, host, port, dbname)
	}

	db, err := sql.Open("mysql", dsn)
	check(err)

	err = db.Ping()
	check(err)

	return db
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
