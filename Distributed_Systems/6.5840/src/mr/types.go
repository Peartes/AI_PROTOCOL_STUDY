package mr

type JobType string

var Map, Reduce JobType = "Map", "Reduce"

type MapJob struct {
	JobId    int    // the job id
	SplitFile string // the split file this job is working on
	NReduce   int    // number of intermediate files to split intermediate results into
	TimeStarted string // the time this job started
}

type ReduceJob[T any] struct {
	JobId int // the job id; 
	IntermediateFilePointer T // the partition key for this job to handle 
	TimeStarted string // the time this job started
}

type Job[T comparable] struct {
	MapJob
	ReduceJob[T] 
}

type IntermediateFile[T comparable] struct {
	PartitionId T
	FileName string
}
type MapJobIntermediateFiles[T comparable] struct {
	IntermediateFiles []IntermediateFile[T] // the intermediate files where all the intermediate k/v pairs are stored based on the partition algo
}

type JobDoneReq[T comparable] struct {
	JobType JobType
	Job Job[T]
	MapJobPartitions MapJobIntermediateFiles[T] // all the files where the intermediate map results are stored
	Err string // was there an error ?
}

type JobDoneReply struct {}

type GetJobRequest struct {}

type GetJobReply[T comparable] struct {
	JobType   JobType // "Map" or "Reduce"
	Job Job[T]
	Files       []string // input files for Map tasks, or intermediate files for Reduce tasks
	Wait bool // flag to instruct the worker to wait some time for a job assignment
	Exit bool // flag to instruct the worker to exit
}

// for sorting by key.
type ByKey []KeyValue

// for sorting by key.
func (a ByKey) Len() int           { return len(a) }
func (a ByKey) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByKey) Less(i, j int) bool { return a[i].Key < a[j].Key }