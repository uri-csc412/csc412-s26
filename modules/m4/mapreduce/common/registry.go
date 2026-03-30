package common

import (
	"m4/mapreduce"
	"m4/mapreduceapps"
)

// JobFns stores a pair of map/reduce functions
type JobFns struct {
	Map    func(string, string) []mapreduce.KeyValue
	Reduce func(string, []string) string
}

// Registry is a global lookup for all supported jobs
var Registry = map[string]JobFns{
	"wc": {
		Map:    mapreduceapps.MapWC,
		Reduce: mapreduceapps.ReduceWC,
	},
	"indexer": {
		Map:    mapreduceapps.MapIndexer,
		Reduce: mapreduceapps.ReduceIndexer,
	},
	"mtiming": {
		Map:    mapreduceapps.MapMTiming,
		Reduce: mapreduceapps.ReduceMTiming,
	},
	"rtiming": {
		Map:    mapreduceapps.MapRTiming,
		Reduce: mapreduceapps.ReduceRTiming,
	},
	"jobcount": {
		Map:    mapreduceapps.MapJobCount,
		Reduce: mapreduceapps.ReduceJobCount,
	},
	"early_exit": {
		Map:    mapreduceapps.MapEarlyExit,
		Reduce: mapreduceapps.ReduceEarlyExit,
	},
	"crash": {
		Map:    mapreduceapps.MapCrash,
		Reduce: mapreduceapps.ReduceCrash,
	},
	"nocrash": {
		Map:    mapreduceapps.MapNoCrash,
		Reduce: mapreduceapps.ReduceNoCrash,
	},
}
