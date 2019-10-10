package main

import (
	"fmt"
	"os"

	"myapp/mig"
)

func main() {
	if len(os.Args) == 1 {
		usage()
		return
	}

	mig.Init()

	switch os.Args[1] {
	case "new":
		mig.New(os.Args)

	case "up", "run":
		mig.Up()

	case "down":
		mig.Down()

	case "log":
		mig.Log()

	case "help":
		usage()

	default:
		fmt.Println("Unknown command:", os.Args[1])
	}
}

func usage() {
	fmt.Println("usage: mig new|up|down|log|help")
}
