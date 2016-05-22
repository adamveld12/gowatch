package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

var (
	shouldErr  = flag.Bool("err", false, "Will exit with a panic after 5 seconds")
	shouldHalt = flag.Bool("halt", false, "Will exit with zero exit code after 5 seconds")
)

func main() {
	flag.Parse()

	if *shouldErr {
		go err()
	} else if *shouldHalt {
		go halt()
	}

	for i := 1; ; i++ {
		time.Sleep(1 * time.Second)

		fmt.Printf("%vs\n", i)
	}
}

func err() {
	<-time.After(5 * time.Second)
	panic("exiting with a failure")
}

func halt() {
	<-time.After(5 * time.Second)
	fmt.Println("exiting successfully")
	os.Exit(0)
}
