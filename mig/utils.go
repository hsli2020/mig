package mig

import (
	"log"
	"os"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func getenv(name, defaultval string) string {
	val := os.Getenv(name)
	if val == "" {
		return defaultval
	}
	return val
}
