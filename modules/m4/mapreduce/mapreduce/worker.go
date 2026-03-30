package mapreduce

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"log"
	"net/rpc"
	"os"
	"sort"
	"time"
)

// ======================
// Types & helpers
// ======================

// Map functions return a slice of KeyValue.
type KeyValue struct {
	Key   string
	Value string
}

// for sorting by key
type ByKey []KeyValue

func (a ByKey) Len() int           { return len(a) }
func (a ByKey) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByKey) Less(i, j int) bool { return a[i].Key < a[j].Key }

// data structure to store key-value pairs for each bucket
type BucketMap map[int][]KeyValue

// use ihash(key) % NReduce to choose the reduce
// task number for each KeyValue emitted by Map.
func ihash(key string) int {
	h := fnv.New32a()
	_, _ = h.Write([]byte(key))
	return int(h.Sum32() & 0x7fffffff)
}

// ======================
// Worker entry point
// ======================

func Worker(
	mapf func(string, string) []KeyValue,
	reducef func(string, []string) string,
) {
	for {
		var reply JobReply
		if !CallGetJob(&reply) {
			// coordinator is gone or unreachable — exit cleanly
			log.Printf("WORKER: coordinator unreachable; exiting.")
			return
		}

		switch reply.JobType {
		case MAP_TASK:
			if err := doMap(reply, mapf); err != nil {
				log.Printf("WORKER: map job %d failed: %v", reply.JobId, err)
				CallNotifyDone(reply.JobId, UNASSIGNED)
			} else {
				// report completion after successful writes
				CallNotifyDone(reply.JobId, FINISHED)
			}

		case REDUCE_TASK:
			if err := doReduce(reply, reducef); err != nil {
				log.Printf("WORKER: reduce job %d failed: %v", reply.JobId, err)
				CallNotifyDone(reply.JobId, UNASSIGNED)
			} else {
				// report completion after atomic rename
				CallNotifyDone(reply.JobId, FINISHED)
			}

		case UNKNOWN_TASK:
			// no work right now — back off briefly
			time.Sleep(500 * time.Millisecond)
		}
	}
}

// ========================
// Map / Reduce processing
// ========================

func doMap(reply JobReply, mapf func(string, string) []KeyValue) error {
	// TODO: Implement the map side of the worker

	// Open input file

	// Read content (slice if DataLength>0, else whole file)
	// var content []byte // use this if u wish ;)

	// populate the key value store using the map function -> mapf

	// Partition the KeyValue pairs into reply.NReduce buckets using ihash(key) % reply.NReduce.
	buckets := make(BucketMap, reply.NReduce) // use this

	// Emit JSON files: mapreduce-M-R (atomic: write .tmp then rename)
	for r := 0; r < reply.NReduce; r++ {
		var buf bytes.Buffer
		if err := json.NewEncoder(&buf).Encode(buckets[r]); err != nil {
			return fmt.Errorf("encode r=%d: %w", r, err)
		}

		tmp := fmt.Sprintf("mapreduce-%d-%d.tmp", reply.JobId, r)
		final := fmt.Sprintf("mapreduce-%d-%d", reply.JobId, r)

		if err := os.WriteFile(tmp, buf.Bytes(), 0644); err != nil {
			return fmt.Errorf("write tmp: %w", err)
		}
		if err := os.Rename(tmp, final); err != nil {
			return fmt.Errorf("rename tmp→final: %w", err)
		}
	}

	// Return an error if any step fails so the caller can decide how to notify the coordinator.
	return nil
}

func doReduce(reply JobReply, reducef func(string, []string) string) error {
	// TODO: Implement the reduce side of the worker

	// TODO: Gather all intermediate files for this reduce -> reply.IntermediateFiles
	// Per file:
	// 	Read the file (e.g., using os.ReadFile)
	// 	Unmarshal the JSON into a temporary []KeyValue
	//  	hint: use json.Unmarshal into a ByKey or []KeyValue)
	//  Append kvPairs pairs into `kvPairs`
	var kvPairs ByKey // use this you beautiful person!

	// TODO: If `len(kvPairs) == 0`, decide how you want to handle no data

	// Sorting here is part of the Reduce phase
	// the shuffle already routed all KV pairs for this bucket to this reducer
	// Now we sort locally so identical keys are adjacent for grouping
	sort.Sort(kvPairs)

	// Write to temp, then atomically rename to final
	// no TODO here :)
	tmp, err := os.CreateTemp("", "mr-out-")
	if err != nil {
		return fmt.Errorf("create temp out: %w", err)
	}
	defer tmp.Close()

	// TODO!
	//  Walk through `all` in order.
	//  Group values with the same key into a []string.
	//  For each key, call reducef(key, values) and write "key value\n"
	//  to the output file (hint: fmt.Fprintf).

	// atomically replace any existing output with the completed temp file
	final := fmt.Sprintf("mapreduce-out-%d", reply.JobId)
	if err := os.Rename(tmp.Name(), final); err != nil {
		return fmt.Errorf("rename reduce out: %w", err)
	}

	// If there are no intermediate files, it is fine to produce an empty output file
	// or no file at all, depending on your design and tests.
	return nil
}

// ==============================
// RPC requests -- do not touch
// ==============================

// send an RPC request to the coordinator; wait for the response,
// usually returns true, return false on any error
//
// IMPORTANT: never log.Fatal here — workers must be able to exit cleanly
func call(rpcname string, args interface{}, reply interface{}) bool {
	sockname := coordinatorSock()
	c, err := rpc.DialHTTP("unix", sockname)
	if err != nil {
		log.Printf("WORKER: dial %s failed: %v", sockname, err)
		return false
	}
	defer c.Close()

	err = c.Call(rpcname, args, reply)
	if err == nil {
		return true
	}

	fmt.Println(err)
	return false
}

// Ask the coordinator for a map/reduce job
// Returns false if coordinator is unreachable or has exited
func CallGetJob(reply *JobReply) bool {
	arg := IntArg{os.Getpid()}
	ok := call("Coordinator.AssignJob", &arg, reply)
	if !ok {
		return false
	}
	if reply.JobId != 0 {
		// Note: your original log printed "Filesize" but used FileOffset
		// Keeping the fields as-is to avoid changing your output format
		fmt.Printf("WORKER: %v Assigned File: %s, Filesize: %v Type: %d\n",
			reply.JobId, reply.FileLocation, reply.FileOffset, reply.JobType)
	}
	return true
}

// notify the coordinator that the job is done
func CallNotifyDone(jobID int, status int) int {
	args := NotifyDoneArgs{JobId: jobID, Status: status}
	var reply IntReply
	if !call("Coordinator.WorkerDone", &args, &reply) {
		return -1
	}
	return reply.Status
}
