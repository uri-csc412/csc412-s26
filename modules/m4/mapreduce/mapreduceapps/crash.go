package mapreduceapps

//
// a MapReduce pseudo-application that sometimes crashes,
// and sometimes takes a long time,
// to test MapReduce's ability to recover.
//

import (
	"m4/mapreduce"
	"sort"
	"strconv"
	"strings"
)

func MapCrash(filename string, contents string) []mapreduce.KeyValue {
	maybeCrashOrSleep()

	kva := []mapreduce.KeyValue{}
	kva = append(kva, mapreduce.KeyValue{Key: "a", Value: filename})
	kva = append(kva, mapreduce.KeyValue{Key: "b", Value: strconv.Itoa(len(filename))})
	kva = append(kva, mapreduce.KeyValue{Key: "c", Value: strconv.Itoa(len(contents))})
	kva = append(kva, mapreduce.KeyValue{Key: "d", Value: "abcdefgxyz"})
	return kva
}

func ReduceCrash(key string, values []string) string {
	maybeCrashOrSleep()

	// sort values to ensure deterministic output.
	vv := make([]string, len(values))
	copy(vv, values)
	sort.Strings(vv)

	val := strings.Join(vv, " ")
	return val
}
