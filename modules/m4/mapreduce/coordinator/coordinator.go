package main

//
// start the coordinator process, which is implemented
// in ../mapreduce/coordinator.go
//
// go run coordinator.go pg*.txt
//
// Please do not change this file.
//

import (
	"fmt"
	"m4/mapreduce"
	"os"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: mapreducecoordinator inputfiles...\n")
		os.Exit(1)
	}

	numberOfJobs := 5
	m := mapreduce.MakeCoordinator(os.Args[1:], numberOfJobs)
	for !m.Done() {
		time.Sleep(time.Second)
	}

	time.Sleep(time.Second)
}
