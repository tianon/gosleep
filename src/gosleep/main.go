package main

import (
	"os"
	"time"
)

func main() {
	d, err := time.ParseDuration(os.Args[1])
	if err != nil {
		panic(err)
	}

	time.Sleep(d)
}
