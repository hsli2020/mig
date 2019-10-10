package mig

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"myapp/mig/database"
)

var db *sql.DB

var migdir string
var migtab string
var dbname string

func Init() {
	migdir = getenv("MIGDIR", "migrations")
	migtab = getenv("MIGTAB", "migration_log")
	dbname = getenv("DBNAME", "bte")

	db = database.Open()

	sql := "CREATE TABLE IF NOT EXISTS " + migtab + " (id varchar(20) not null primary key, action varchar(80) default '', run_at timestamp not null default CURRENT_TIMESTAMP) engine=InnoDB"
	db.Exec(sql)
}

func New(args []string) {
	if len(args) < 3 {
		fmt.Println("usage: mig new <NAME>")
		return
	}

	message := args[2]
	version := time.Now().Format("20060102-150405")

	// create migration up script
	filename := filepath.Join(migdir, version+"_"+message+".up.sql")
	file, err := os.Create(filename)
	check(err)
	file.WriteString("USE " + dbname + ";\n")
	file.Close()
	fmt.Println("Created:", filename)

	// create migration down script
	filename = filepath.Join(migdir, version+"_"+message+".down.sql")
	file, err = os.Create(filename)
	check(err)
	file.WriteString("USE " + dbname + ";\n")
	file.Close()
	fmt.Println("Created:", filename)
}

func Up() {
	pattern := filepath.Join(migdir, "*.up.sql")

	files, err := filepath.Glob(pattern)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, file := range files {
		// load migration script
		script, err := ioutil.ReadFile(file)
		if err != nil {
			continue // next migration script
		}

		// run migration script
		_, err = db.Exec(string(script))
		if err != nil {
			continue // next migration script
		}

		// log the migration for tracking
		basename := filepath.Base(file)
		version := basename[0:15]
		message := basename[16:]

		sql := fmt.Sprintf("INSERT INTO %s set id='%s', action='%s'", migtab, version, message)
		db.Exec(sql)

		fmt.Println("Running:", file)
	}
}

func Down() {
	var version, message, sql string

	// fetch last migration
	sql = fmt.Sprintf("SELECT id, action FROM %s ORDER BY id DESC LIMIT 1", migtab)
	err := db.QueryRow(sql).Scan(&version, &message)
	check(err)

	// find the migration down script
	filename := filepath.Join(migdir, version+"_"+message[0:len(message)-7]+".down.sql")
	script, err := ioutil.ReadFile(filename)
	check(err)

	// run the migration script
	_, err = db.Exec(string(script))
	check(err)

	sql = fmt.Sprintf("DELETE FROM %s WHERE id='%s'", migtab, version)
	_, err = db.Exec(sql)
	check(err)

	fmt.Println("Running:", filename)
}

func Log() {
	var version, message, sql string

	sql = fmt.Sprintf("SELECT id, action FROM %s ORDER BY id", migtab)
	rows, err := db.Query(sql)
	check(err)
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&version, &message)
		check(err)
		fmt.Println(version, message)
	}
}
