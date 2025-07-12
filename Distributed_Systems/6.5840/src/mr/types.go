package mr

type TaskType string

var Map, Reduce TaskType

type MapJobPartitions struct {
	PartitionId int // the file where the map job intermediate res is stored
	Path string // where is this file located ? but more like what it's name is
}

type JobDoneReq struct {
	TaskNumber int // the map or reduce job id
	TaskType TaskType
	MapJobPartitions []MapJobPartitions // all the files where the intermediate map results are stored
	err error // was there an error ?
}

type JobDoneReply struct {

}

type GetTaskRequest struct {
}

type GetTaskReply struct {
	TaskType   TaskType // "Map" or "Reduce"
	TaskNumber int    // task number for this Map or Reduce task
	Files       []string // input files for Map tasks, or intermediate files for Reduce tasks
	Partitions int // number of intermediate files to create for a map task
}

// for sorting by key.
type ByKey []KeyValue

// for sorting by key.
func (a ByKey) Len() int           { return len(a) }
func (a ByKey) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByKey) Less(i, j int) bool { return a[i].Key < a[j].Key }