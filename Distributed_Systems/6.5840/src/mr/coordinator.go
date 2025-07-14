package mr

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"sync"
	"time"
)

/*
Coordinator is the main struct that holds the state of the MapReduce job.
All file splits to be processed are stored in splits.
The pendingMapJobs are the map jobs not yet processed. when all map jobs are done, this array should be empty
The processingMapJobs is an array of map jobs that are currently being processed.
The completedMapJobs is an array of all completed map jobs ... when all map jobs are done this array's length is len(splits)
# intermediateFiles is the sorted list of intermediate files produced by map jobs workers.
# This list is sorted to group together all similar intermediate files that must contain similar intermediate keys
# The key is the intermediate file partition and the value is all intermediate files fileName (i.e. files created to hold all intermediate k/v that fall into a partition)
# This works because there are at most R intermediate files (from the partitioning algo)
The pendingReduces is an array of all pending reduce jobs. When all reduce jobs are done, this array must be empty
The processingReduces is the array of currently processing reduce jobs
The completedReduces is an array of completed reduce jobs. When all reduce jobs are done this array's length should be len(keys(intermediateFile))
# since there are at most R (number of reduce jobs) intermediate files
nReduce is the number of reduce jobs to create
nMap is the number of map jobs to create
*/
type Coordinator[T comparable] struct {
	// Your definitions here.
	splits               []string
	pendingMapJobs       []MapJob
	processingMapJobs    []MapJob
	completedMapJobs     []MapJob
	intermediateFiles    map[T][]string
	pendingReduceJobs    []ReduceJob[T]
	processingReduceJobs []ReduceJob[T]
	completedReduceJobs  []ReduceJob[T]
	nMap                 int
	nReduce              int
	done                 chan struct{} // channel to signal when the coordinator is done
}

var mu sync.Mutex

// we want to keep track of how long a task has been running
// so we can detect if a worker has crashed. we use 10s
// as the default timeout.
func checkAllTaskStatus[T comparable](c *Coordinator[T]) {
	mu.Lock()
	defer mu.Unlock()
	for _, job := range c.processingMapJobs {
		timeStarted, err := time.Parse(time.RFC3339, job.TimeStarted)
		if err != nil {
			log.Printf("could not parse time started for job %d: %v", job.JobId, err)
			continue
		}
		// if the job has been running for more than 10 seconds, we can assume
		// that the worker has crashed and we can re-assign the job
		if time.Since(timeStarted) > 10*time.Second {
			log.Printf("Map job %d has been running for more than 10 seconds, it might have crashed", job.JobId)
			// we can re-assign this job to another worker
			c.pendingMapJobs = append(c.pendingMapJobs, job)
			jobIdx := findMapJobAtIndex(c.processingMapJobs, job.JobId)
			if jobIdx == -1 {
				log.Printf("Map job %d is not in processing jobs, it might have been re-assigned already", job.JobId)
				continue
			}
			// remove the job from processing jobs
			c.processingMapJobs = removeMapJobAtIndex(c.processingMapJobs, jobIdx)
		}
	}
	for _, job := range c.processingReduceJobs {
		timeStarted, err := time.Parse(time.RFC3339, job.TimeStarted)
		if err != nil {
			log.Printf("could not parse time started for job %d: %v", job.JobId, err)
			continue
		}
		// if the job has been running for more than 10 seconds, we can assume
		// that the worker has crashed and we can re-assign the job
		if time.Since(timeStarted) > 10*time.Second {
			log.Printf("Reduce job %d has been running for more than 10 seconds, it might have crashed", job.JobId)
			// we can re-assign this job to another worker
			c.pendingReduceJobs = append(c.pendingReduceJobs, job)
			jobIdx := findMapJobAtIndex(c.processingMapJobs, job.JobId)
			if jobIdx == -1 {
				log.Printf("Map job %d is not in processing jobs, it might have been re-assigned already", job.JobId)
				continue
			}
			c.processingReduceJobs = removeReduceJobAtIndex(c.processingReduceJobs, jobIdx)
		}
	}
}

// Your code here -- RPC handlers for the worker to call.

// find the index of a job in the job array
func findMapJobAtIndex(s []MapJob, key int) int {
	for idx, job := range s {
		if job.JobId == key {
			return idx
		}
	}
	return -1
}
func removeMapJobAtIndex(s []MapJob, i int) []MapJob {
	return append(s[:i], s[i+1:]...)
}

// find the index of a job in the job array
func findReduceJobAtIndex[T comparable](s []ReduceJob[T], key int) int {
	for idx, job := range s {
		if job.JobId == key {
			return idx
		}
	}
	return -1
}
func removeReduceJobAtIndex[T comparable](s []ReduceJob[T], i int) []ReduceJob[T] {
	return append(s[:i], s[i+1:]...)
}

func (c *Coordinator[T]) JobDone(args *JobDoneReq[T], reply *JobDoneReply) error {
	fmt.Println("JobDone called")
	mu.Lock()
	defer mu.Unlock()
	if args.JobType == Map {
		fmt.Printf("JobDone called for map job %d\n", args.Job.MapJob.JobId)
		jobIdx := findMapJobAtIndex(c.processingMapJobs, args.Job.MapJob.JobId)
		if jobIdx == -1 {
			// this finished job is not is not supposed to be processing
			// might be a stale worker whose job has been re-assigned
			fmt.Printf("JobDone called for map job %d but it is not in processing jobs\n", args.Job.MapJob.JobId)
			return nil
		}
		if args.Err == "" {
			// this map job is complete
			c.completedMapJobs = append(c.completedMapJobs, args.Job.MapJob)
			// sort the response of the map jobs which is the intermediate files using their partition id
			for _, partitions := range args.MapJobPartitions.IntermediateFiles {
				c.intermediateFiles[partitions.PartitionId] = append(c.intermediateFiles[partitions.PartitionId], partitions.FileName)
			}
			// if this is the last map job, let's assign the reduce jobs
			if len(c.completedMapJobs) == len(c.splits) {
				i := 0
				for partitionId, _ := range c.intermediateFiles {
					c.pendingReduceJobs = append(c.pendingReduceJobs, ReduceJob[T]{JobId: i, IntermediateFilePointer: partitionId})
					i = i + 1
				}
			}
		} else {
			// an error occurred in processing so let's add the job back to pending
			fmt.Printf("JobDone called for map job %d but it had an error: %v\n", args.Job.MapJob.JobId, args.Err)
			c.pendingMapJobs = append(c.pendingMapJobs, args.Job.MapJob)
		}
		// remove from processing
		c.processingMapJobs = removeMapJobAtIndex(c.processingMapJobs, jobIdx)
	} else {
		fmt.Printf("JobDone called for reduce job %d\n", args.Job.ReduceJob.JobId)
		jobIdx := findReduceJobAtIndex(c.processingReduceJobs, args.Job.ReduceJob.JobId)
		if jobIdx == -1 {
			// this must be a stale worker response
			fmt.Printf("JobDone called for reduce job %d but it is not in processing jobs\n", args.Job.ReduceJob.JobId)
			return nil
		}
		if args.Err == "" {
			// this reduce job completed successfully
			c.completedReduceJobs = append(c.completedReduceJobs, args.Job.ReduceJob)
		} else {
			fmt.Printf("JobDone called for reduce job %d but it had an error: %v\n", args.Job.ReduceJob.JobId, args.Err)
			c.pendingReduceJobs = append(c.pendingReduceJobs, args.Job.ReduceJob)
		}
		// remove the job from processing jobs
		c.processingReduceJobs = removeReduceJobAtIndex(c.processingReduceJobs, jobIdx)
	}
	return nil
}

// start a thread that listens for RPCs from worker.go
func (c *Coordinator[T]) server() {
	rpc.RegisterName("Coordinator", c)
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

// main/mrcoordinator.go calls Done() periodically to find out
// if the entire job has finished.
func (c *Coordinator[T]) Done() bool {
	ret := false

	// the map reduce library is done when all pending and processing reduce and map tasks are done
	// and the number of the completed reduce tasks is equal to the total number of reduce jobs
	if len(c.pendingReduceJobs) == 0 && len(c.processingReduceJobs) == 0 && len(c.completedReduceJobs) > 0 && len(c.completedReduceJobs) == len(c.intermediateFiles) {
		ret = true
		close(c.done)
	}

	return ret
}

func (c *Coordinator[T]) RequestJob(args *GetJobRequest, reply *GetJobReply[T]) error {
	// we only reduce jobs after all map jobs are done
	fmt.Println("RequestJob called")
	mu.Lock()
	defer mu.Unlock()
	if len(c.pendingMapJobs) > 0 {
		// there are still map jobs
		job := c.pendingMapJobs[0]
		// remove this pending job
		c.pendingMapJobs = c.pendingMapJobs[1:]
		reply.JobType = Map
		reply.Job = Job[T]{
			MapJob:    job,
			ReduceJob: ReduceJob[T]{},
		}
		reply.Files = []string{job.SplitFile}
		// move the job into the processing array
		job.TimeStarted = time.Now().Format(time.RFC3339)
		c.processingMapJobs = append(c.processingMapJobs, job)
		fmt.Printf("Worker assigned map job %d for file %s\n", job.JobId, job.SplitFile)
	} else if len(c.processingMapJobs) > 0 {
		// there are still running map jobs
		// we will wait until they're done
		reply.Wait = true
		fmt.Println("Worker waiting for map jobs to complete")
	} else {
		// all map jobs must be done
		if len(c.completedMapJobs) < len(c.splits) {
			// TODO: this does not have to be a panic since we can just run the missing job
			// but it shows some problem in our logic most likely race condition
			panic("some map jobs are missing")
		}
		// if all reduce jobs are done, tell the worker to exit
		if len(c.pendingReduceJobs) == 0 && len(c.processingReduceJobs) == 0 && len(c.completedReduceJobs) > 0 && len(c.completedReduceJobs) == len(c.intermediateFiles) {
			reply.Exit = true
			fmt.Println("All jobs are done, worker should exit")
			return nil
		}
		// assign a reduce job
		if len(c.pendingReduceJobs) == 0 {
			// no pending reduce jobs, but there are processing reduce jobs
			reply.Wait = true
			fmt.Println("Worker waiting for reduce jobs to complete")
			return nil
		}
		// there are pending reduce jobs
		job := c.pendingReduceJobs[0]
		// remove this pending job
		c.pendingReduceJobs = c.pendingReduceJobs[1:]
		reply.JobType = Reduce
		reply.Job = Job[T]{
			MapJob{},
			job,
		}
		reply.Files = c.intermediateFiles[job.IntermediateFilePointer]
		// move the job into the processing array
		job.TimeStarted = time.Now().Format(time.RFC3339)
		c.processingReduceJobs = append(c.processingReduceJobs, job)
		fmt.Printf("Worker assigned reduce job %d for intermediate files %v\n", job.JobId, reply.Files)
	}
	return nil
}

// create a Coordinator.
// main/mrcoordinator.go calls this function.
// nReduce is the number of reduce tasks to use.
func MakeCoordinator[T comparable](files []string, nReduce int) *Coordinator[T] {
	c := Coordinator[T]{
		splits:               files,
		processingMapJobs:    []MapJob{},
		completedMapJobs:     []MapJob{},
		pendingReduceJobs:    []ReduceJob[T]{},
		processingReduceJobs: []ReduceJob[T]{},
		completedReduceJobs:  []ReduceJob[T]{},
		intermediateFiles:    map[T][]string{},
		nMap:                 len(files),
		nReduce:              nReduce,
	}

	// build the pending map jobs
	for i, split := range files {
		c.pendingMapJobs = append(c.pendingMapJobs, MapJob{JobId: i, SplitFile: split, NReduce: nReduce})
	}

	c.server()
	// start a thread that checks the status of all tasks
	go func() {
		for {
			select {
			case <-c.done:
				fmt.Println("stale checker shutting down")
				return
			default:
				checkAllTaskStatus(&c)
				time.Sleep(10 * time.Second)
			}
		}
	}()
	return &c
}
