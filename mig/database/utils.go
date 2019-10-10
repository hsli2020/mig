package database

import (
	"log"
	"os"
)

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
