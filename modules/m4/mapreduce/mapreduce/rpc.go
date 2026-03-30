package mapreduce

//
// RPC definitions
//

import (
	"fmt"
	"os"
	"runtime"
)

// ======================
// Types & helpers
// ======================

const (
	UNASSIGNED = -1
	FINISHED   = 0
	ASSIGNED   = 1
)

const (
	UNKNOWN_TASK = 0
	MAP_TASK     = 1
	REDUCE_TASK  = 2
)

// A simple Integer argument
type IntArg struct {
	Status int
}

// A simple Integer reply, used for confirmation
type IntReply struct {
	Status int
}

// The Coordinate replies with Map Job request
// and gives details on the job to perform.
type JobReply struct {
	JobId             int
	JobType           int
	FileLocation      string
	FileOffset        int64
	DataLength        int64
	NReduce           int
	IntermediateFiles []string
}

// The worker must notify the Coordinator
// that the map/reduce is done and where the
// intermediate files are located.
type NotifyDoneArgs struct {
	JobId    int
	Status   int
	Location string
}

// ======================
// Socket
// ======================

// the coordinator will listen on socket/412-mapreduce.sock
// and every worker will try to dial exactly that
func coordinatorSock() string {
	// project-local socket dir for testing
	socketDir := "socket"

	// make sure it exists (coordinator will create; workers calling this is harmless)
	_ = os.MkdirAll(socketDir, 0755)

	// everyone (coordinator + workers) uses the SAME name
	if runtime.GOOS == "windows" {
		// still give Windows something, even though unix sockets aren't used the same way
		return fmt.Sprintf("%s%c412-mapreduce.sock", socketDir, os.PathSeparator)
	}

	return fmt.Sprintf("%s%c412-mapreduce.sock", socketDir, os.PathSeparator)
}
