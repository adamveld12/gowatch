// test is a simple app that counts the number of seconds its been running.
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
		fmt.Print(fmtIteration(i))
	}
}

func fmtIteration(i int) string {
	pluralized := ""

	if i > 1 || i < 1 {
		pluralized = "s"
	}

	return fmt.Sprintf("%v second%v and counting!\n", i, pluralized)
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
