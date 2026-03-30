package mapreduceapps

//
// a MapReduce pseudo-application to test that workers
// execute reduce tasks in parallel.
//

import (
	"fmt"
	"m4/mapreduce"
	"os"
	"syscall"
	"time"
)

func nparallelR(phase string) int {
	// create a file so that other workers will see that
	// we're running at the same time as them.
	pid := os.Getpid()
	myfilename := fmt.Sprintf("mapreduce-worker-%s-%d", phase, pid)
	err := os.WriteFile(myfilename, []byte("x"), 0666)
	if err != nil {
		panic(err)
	}

	// are any other workers running?
	// find their PIDs by scanning directory for mapreduce-worker-XXX files.
	dd, err := os.Open(".")
	if err != nil {
		panic(err)
	}
	names, err := dd.Readdirnames(1000000)
	if err != nil {
		panic(err)
	}
	ret := 0
	for _, name := range names {
		var xpid int
		pat := fmt.Sprintf("mapreduce-worker-%s-%%d", phase)
		n, err := fmt.Sscanf(name, pat, &xpid)
		if n == 1 && err == nil {
			err := syscall.Kill(xpid, 0)
			if err == nil {
				// if err == nil, xpid is alive.
				ret += 1
			}
		}
	}
	dd.Close()

	time.Sleep(1 * time.Second)

	err = os.Remove(myfilename)
	if err != nil {
		panic(err)
	}

	return ret
}

func MapRTiming(filename string, contents string) []mapreduce.KeyValue {

	kva := []mapreduce.KeyValue{}
	kva = append(kva, mapreduce.KeyValue{Key: "a", Value: "1"})
	kva = append(kva, mapreduce.KeyValue{Key: "b", Value: "1"})
	kva = append(kva, mapreduce.KeyValue{Key: "c", Value: "1"})
	kva = append(kva, mapreduce.KeyValue{Key: "d", Value: "1"})
	kva = append(kva, mapreduce.KeyValue{Key: "e", Value: "1"})
	kva = append(kva, mapreduce.KeyValue{Key: "f", Value: "1"})
	kva = append(kva, mapreduce.KeyValue{Key: "g", Value: "1"})
	kva = append(kva, mapreduce.KeyValue{Key: "h", Value: "1"})
	kva = append(kva, mapreduce.KeyValue{Key: "i", Value: "1"})
	kva = append(kva, mapreduce.KeyValue{Key: "j", Value: "1"})
	return kva
}

func ReduceRTiming(key string, values []string) string {
	n := nparallelR("reduce")

	val := fmt.Sprintf("%d", n)

	return val
}
