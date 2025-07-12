package mr

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"sync"
)


type Coordinator struct {
	// Your definitions here.
	pendingSplits []string // the original file splits (fileNames) not yet worked on
	processingMaps map[int]string // a map of job id to the split they're working
	intermediateFiles map[int][]string // map of the reduce jobs (partition) not yet worked on to their intermediate file
	pendingReduces []int // the reduces not yet done
	processingReduce  map[int]int // map of reduce job to the current reduce worker
}

var mu sync.Mutex
// we want to keep track of how long a task has been running
// so we can detect if a worker has crashed. we use 10s
// as the default timeout.
func checkAllTaskStatus() {}
// Your code here -- RPC handlers for the worker to call.

func (c *Coordinator) JobDone(args *JobDoneReq , reply *JobDoneReply) error {
	if args.TaskType == Map {
		// remove from processing
		mu.Lock()
		if args.err == nil {
			delete(c.processingMaps, args.TaskNumber)
			// set the response of this map job as pending reduces jobs
			for _, reduceJobId := range args.MapJobPartitions {
				c.intermediateFiles[reduceJobId.PartitionId] = append(c.intermediateFiles[reduceJobId.PartitionId], reduceJobId.Path)
				c.pendingReduces = append(c.pendingReduces, reduceJobId.PartitionId)
			}
		} else {
			// an error occurred in processing so let's add the job back to pending
			failedSplitFile := c.processingMaps[args.TaskNumber]
			c.pendingSplits = append(c.pendingSplits, failedSplitFile)
		}
		mu.Unlock()
	} else {
		// remove this task from the processing reduces
		mu.Lock()
		if args.err == nil {
			delete(c.processingReduce, args.TaskNumber)
		} else {
			c.pendingReduces = append(c.pendingReduces, args.TaskNumber)
		}
		mu.Unlock()
	}
	return nil
}

//
// an example RPC handler.
//
// the RPC argument and reply types are defined in rpc.go.
//
func (c *Coordinator) Example(args *ExampleArgs, reply *ExampleReply) error {
	reply.Y = args.X + 1
	return nil
}


//
// start a thread that listens for RPCs from worker.go
//
func (c *Coordinator) server() {
	rpc.Register(c)
	rpc.HandleHTTP()
	//l, e := net.Listen("tcp", ":1234")
	sockname := coordinatorSock()
	os.Remove(sockname)
	l, e := net.Listen("unix", sockname)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}

//
// main/mrcoordinator.go calls Done() periodically to find out
// if the entire job has finished.
//
func (c *Coordinator) Done() bool {
	ret := false

	// Your code here.
	if (len(c.pendingSplits) == 0 && len(c.processingMaps) == 0 && len(c.pendingReduces)= 0 && c.processingReduce == 0) {
		 ret = true
	}

	return ret
}

//
// create a Coordinator.
// main/mrcoordinator.go calls this function.
// nReduce is the number of reduce tasks to use.
//
func MakeCoordinator(files []string, nReduce int) *Coordinator {
	c := Coordinator{}

	// Your code here.

	c.server()
	return &c
}
