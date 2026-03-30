package mapreduceapps

//
// same as crash.go but doesn't actually crash
//

import (
	"m4/mapreduce"
	"sort"
	"strconv"
	"strings"
)

func MapNoCrash(filename string, contents string) []mapreduce.KeyValue {
	maybeCrash()

	kva := []mapreduce.KeyValue{}
	kva = append(kva, mapreduce.KeyValue{Key: "a", Value: filename})
	kva = append(kva, mapreduce.KeyValue{Key: "b", Value: strconv.Itoa(len(filename))})
	kva = append(kva, mapreduce.KeyValue{Key: "c", Value: strconv.Itoa(len(contents))})
	kva = append(kva, mapreduce.KeyValue{Key: "d", Value: "abcdefg"})
	return kva
}

func ReduceNoCrash(key string, values []string) string {
	maybeCrash()

	// sort values to ensure deterministic output.
	vv := make([]string, len(values))
	copy(vv, values)
	sort.Strings(vv)

	val := strings.Join(vv, " ")
	return val
}
