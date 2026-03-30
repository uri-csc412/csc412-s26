package main

// simple sequential MapReduce
//
// usage:
//   go run ./sequential wc ./data/pg*.txt

import (
	"fmt"
	"io"
	"log"
	"os"
	"sort"

	"m4/common"
	"m4/mapreduce"
)

// for sorting by key.
type ByKey []mapreduce.KeyValue

func (a ByKey) Len() int           { return len(a) }
func (a ByKey) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByKey) Less(i, j int) bool { return a[i].Key < a[j].Key }

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "usage: mapreducesequential <jobname> <inputfiles...>\n")
		os.Exit(1)
	}

	jobname := os.Args[1]
	fns, ok := common.Registry[jobname]
	if !ok {
		log.Fatalf("unknown job %q", jobname)
	}

	mapf := fns.Map
	reducef := fns.Reduce

	var intermediate []mapreduce.KeyValue

	for _, filename := range os.Args[2:] {
		f, err := os.Open(filename)
		if err != nil {
			log.Fatalf("cannot open %v: %v", filename, err)
		}
		data, err := io.ReadAll(f)
		if err != nil {
			log.Fatalf("cannot read %v: %v", filename, err)
		}
		_ = f.Close()

		kva := mapf(filename, string(data))
		intermediate = append(intermediate, kva...)
	}

	sort.Sort(ByKey(intermediate))

	outName := "mapreduce-out-0"
	out, err := os.Create(outName)
	if err != nil {
		log.Fatalf("cannot create %v: %v", outName, err)
	}
	defer out.Close()

	i := 0
	for i < len(intermediate) {
		j := i + 1
		for j < len(intermediate) && intermediate[j].Key == intermediate[i].Key {
			j++
		}
		var values []string
		for k := i; k < j; k++ {
			values = append(values, intermediate[k].Value)
		}

		output := reducef(intermediate[i].Key, values)
		fmt.Fprintf(out, "%v %v\n", intermediate[i].Key, output)

		i = j
	}
}
