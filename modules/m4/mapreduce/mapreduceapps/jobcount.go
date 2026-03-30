package mapreduceapps

//
// a MapReduce pseudo-application that counts the number of times map/reduce
// tasks are run, to test whether jobs are assigned multiple times even when
// there is no failure.
//

import (
	"fmt"
	"m4/mapreduce"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

var count int

func MapJobCount(filename string, contents string) []mapreduce.KeyValue {
	me := os.Getpid()
	f := fmt.Sprintf("mapreduce-worker-jobcount-%d-%d", me, count)
	count++
	err := os.WriteFile(f, []byte("x"), 0666)
	if err != nil {
		panic(err)
	}
	time.Sleep(time.Duration(2000+rand.Intn(3000)) * time.Millisecond)
	return []mapreduce.KeyValue{{Key: "a", Value: "x"}}
}

func ReduceJobCount(key string, values []string) string {
	files, err := os.ReadDir(".")
	if err != nil {
		panic(err)
	}
	invocations := 0
	for _, f := range files {
		if strings.HasPrefix(f.Name(), "mapreduce-worker-jobcount") {
			invocations++
		}
	}
	return strconv.Itoa(invocations)
}
