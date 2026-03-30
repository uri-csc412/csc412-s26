package mapreduce

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"regexp"
	"strconv"
	"sync"
	"time"
)

// ======================
// Types & helpers
// ======================

// Coordinator holds all state about map and reduce jobs
// // do not edit ;)
type Coordinator struct {
	mapJobs    []map_job    // Slice of map jobs
	reduceJobs []reduce_job // Slice of reduce jobs
	numReduce  int          // Total number of reduce jobs
	mapDone    int          // 0 = map phase, 1 = reduce phase
	mu         sync.Mutex   // mutex to lock mapJobs and reduceJobs and mapDone
}

// map_job represents a single map task
// // do not edit ;)
type map_job struct {
	job_id        int
	job_type      int
	status        int
	file_location string
	file_index    int64
	m_size        int64
	timeout       int
}

// reduce_job represents a single reduce task
// // do not edit ;)
type reduce_job struct {
	job_id            int
	status            int
	intermediateFiles []string
	timeout           int
}

// ======================
// Coordinator Functions
// ======================

// GetIntermediateFiles scans the current directory for intermediate files
// that belong to a given reduce bucket. We look for files like
//
//	mapreduce-<mapJobID>-<reduceNum>
//
// // likely no need to edit ;)
func GetIntermediateFiles(rNum int) []string {
	var matchedFiles []string
	pattern := fmt.Sprintf(`^mapreduce-(\d+)-%d$`, rNum)
	regex, err := regexp.Compile(pattern)
	if err != nil {
		fmt.Println("Regex error:", err)
		return nil
	}

	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current working directory:", err)
		return nil
	}

	// Read all files in the current directory
	files, err := os.ReadDir(cwd)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return nil
	}

	for _, file := range files {
		if file.IsDir() {
			continue // skip directories
		}
		fileName := file.Name()
		if regex.MatchString(fileName) {
			matchedFiles = append(matchedFiles, fileName)
		}
	}
	return matchedFiles
}

// Start a thread that listens for RPCs from worker.go
// // do not edit ;)
func (c *Coordinator) server() {
	rpc.Register(c)
	rpc.HandleHTTP()

	// make sure our socket dir exists
	_ = os.MkdirAll("socket", 0o755)

	sockname := coordinatorSock()
	// remove any stale socket file
	_ = os.Remove(sockname)

	l, e := net.Listen("unix", sockname)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}

// AssignJob Routine
// assigns a map or reduce job to worker
func (c *Coordinator) AssignJob(proc_id *IntArg, reply *JobReply) error {
	// TODO: what should we lock

	if c.mapDone == 0 {
		// TODOS
		// Assign worker a map job
		// Check Status, Mark this job as assigned AND reset its timeout
		// should we return nil?
	} else {
		// TODOS
		// Assign worker a reduce job
		// Check Status, Mark this job as assigned AND reset its timeout
		// should we return nil?
	}

	// no work right now
	reply.JobType = UNKNOWN_TASK
	return nil
}

// WorkerDone - end condition
// Processes the reply from worker
func (c *Coordinator) WorkerDone(args *NotifyDoneArgs, reply *IntReply) error {
	// TODO: what should we lock

	if c.mapDone == 0 {
		// TODOS
		// Find matching map job
		// Return -1 to worker for an already completed job
		// mark as the status the worker reported (usually FINISHED)
		// should we return nil?
	} else {
		// TODOS
		// Find matching reduce job
		// Return -1 to worker for an already completed job
		// should we return nil?
	}

	// notify worker it's job is done
	reply.Status = 0
	return nil
}

// Done reports whether the current phase is finished.
// CALLER should sleep, not this method.
func (c *Coordinator) Done() bool {
	// TODO: what should we lock

	// set allJobsFinished based on whether all jobs in the current phase are finished
	allJobsFinished := true

	total := 0
	complete := 0

	if c.mapDone == 0 {
		// TODO: MAP PHASE
		// Count total and complete jobs
		// For ASSIGNED jobs, increment timeout and possibly requeue
		// by setting status back to UNASSIGNED when timeout is too large

		// consider being less aggressive about requeueing map tasks:
		// check polling in MakeCoordinator
		// const mapRequeuePolls = 12

		fmt.Printf("COORDINATOR: Status: Total Map Jobs: %d Complete Map Jobs: %d\n", total, complete)

	} else {
		// TODO: REDUCE PHASE
		// Count total and complete jobs
		// For ASSIGNED jobs, increment timeout and possibly requeue

		// reduce tasks can stay more aggressive if you want
		// const reduceRequeuePolls = 2

		fmt.Printf("COORDINATOR: Status: Total Reduce Jobs: %d Complete Reduce Jobs: %d\n", total, complete)
	}

	// TODO: return true only when all jobs in the current phase are FINISHED
	return allJobsFinished
}

// create a Coordinator
// main/mapreducecoordinator.go calls this function
// nReduce is the number of reduce tasks to use
func MakeCoordinator(files []string, nReduce int) *Coordinator {

	c := Coordinator{
		numReduce: nReduce,
	}

	// Parameter values
	map_input_size := int64(64 * 1024 * 1024) // Map input data size == 64MB
	job_id := 1000                            // Starting map job ID

	// Read input data files
	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			fmt.Printf("Cannot access file: %s, error: %v\n", file, err)
			return &c
		}
		fi, err := f.Stat()
		f.Close()
		if err != nil {
			fmt.Printf("Cannot access file: %s, error: %v\n", file, err)
			return &c
		}

		fileSize := fi.Size()
		var offset int64 = 0
		var chunkSize int64 = map_input_size
		// Split input data into "M" pieces of map_input_size
		for offset < fileSize {

			if fileSize-offset < chunkSize {
				chunkSize = fileSize - offset
			}

			// Add new map job to coordinator's map slice
			map_job := map_job{
				job_id:        job_id,
				job_type:      MAP_TASK,
				status:        UNASSIGNED,
				file_location: file,
				file_index:    offset,
				m_size:        chunkSize,
				timeout:       0,
			}
			c.mapJobs = append(c.mapJobs, map_job)
			offset += chunkSize
			job_id = job_id + 1
		}
	}
	fmt.Println("COORDINATOR: ready to assign Map jobs.")

	// Thread that listens for Jobs
	c.server()

	// allow test scripts to speed up polling
	poll := 500 * time.Millisecond
	if v := os.Getenv("POLL_MS"); v != "" {
		if ms, err := strconv.Atoi(v); err == nil && ms > 0 {
			poll = time.Duration(ms) * time.Millisecond
		}
	}

	// Done calls periodically -- MAP PHASE
	for {
		if c.Done() {
			fmt.Println("COORDINATOR: Map jobs are complete. Creating reduce jobs...")
			break
		}
		time.Sleep(poll)
	}

	// switch to reduce phase
	// c.mu.Lock() // hmmm should we be locking here?
	c.mapDone = 1

	// Generate nReduce reduce jobs
	for rNum := 0; rNum < nReduce; rNum++ {
		// grep for the intermediate files of bucket R
		var intermediateFiles = GetIntermediateFiles(rNum)
		reduce_job := reduce_job{rNum, UNASSIGNED, intermediateFiles, 0}
		// Add new reduce job to coordinator's reduce slice
		c.reduceJobs = append(c.reduceJobs, reduce_job)
	}
	// when should we call? c.mu.Unlock() - think about critical sections of code!

	// Done calls periodically -- REDUCE PHASE
	for {
		if c.Done() {
			fmt.Println("COORDINATOR: Reduce jobs are complete. Exiting...")
			break
		}
		time.Sleep(poll)
	}

	return &c
}
