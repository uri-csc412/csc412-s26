package main

//
// start a worker process, which is implemented
// in ../mapreduce/worker.go
//
// go run ./worker wc
//

import (
	"fmt"
	"log"
	"os"

	"m4/common"
	"m4/mapreduce"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: mapreduceworker <jobname>\n")
		os.Exit(1)
	}

	jobname := os.Args[1]
	fns, ok := common.Registry[jobname]
	if !ok {
		log.Fatalf("unknown job %q", jobname)
	}

	// this is your real worker loop in m4/mapreduce
	mapreduce.Worker(fns.Map, fns.Reduce)
}
