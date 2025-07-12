package mr

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"log"
	"net/rpc"
	"os"
	"sort"
	"strconv"
	"time"
)

// Map functions return a slice of KeyValue.
type KeyValue struct {
	Key   string
	Value string
}

// use ihash(key) % NReduce to choose the reduce
// task number for each KeyValue emitted by Map.
func ihash(key string) int {
	h := fnv.New32a()
	h.Write([]byte(key))
	return int(h.Sum32() & 0x7fffffff)
}

// main/mrworker.go calls this function.
func Worker(mapf func(string, string) []KeyValue,
	reducef func(string, []string) string) {

	// Your worker implementation here.

	// workers in a loop, ask the coordinator for work,
	// read the task's input from one or more files,
	// execute the task, write the task's output to one or more files,
	// and again ask the coordinator for a new task
	for {
		// ask the coordinator for a task
		// make an RPC call to the coordinator on the method GetTask
		args := GetJobRequest{}
		var reply GetJobReply[int]
		var doneArgs JobDoneReq[int]
		doneReply := JobDoneReply{}
		ok := call("Coordinator.RequestJob", &args, &reply)
		if !ok {
			return // coordinator is not available, exit the worker
		}
		if reply.Wait {
			time.Sleep(time.Millisecond * 1000)
			continue
		}
		if reply.JobType == Map {
			// read the input files
			file, err := os.Open(reply.Files[0]) // for a map function, we send just one file per map
			if err != nil {
				log.Printf("map worker %d could not open file %s there's no need to continue with map job\n", reply.Job.MapJob.JobId, reply.Files[0])
				doneArgs.err = fmt.Errorf("map worker %d could not open file %s there's no need to continue with map job", reply.Job.MapJob.JobId, reply.Files[0])
				ok := call("Coordinator.JobDone", &doneArgs, &doneReply)
				if !ok {
					// coordinator is done let's exit
					return
				}
				continue
			}
			content, err := ioutil.ReadAll(file)
			if err != nil {
				log.Printf("map worker %d could not read file %s there's no need to continue with map job\n", reply.Job.MapJob.JobId, reply.Files[0])
				doneArgs.err = fmt.Errorf("map worker %d could not read file %s there's no need to continue with map job", reply.Job.MapJob.JobId, reply.Files[0])
				ok = call("Coordinator.JobDone", &doneArgs, &doneReply)
				if !ok {
					// coordinator is done let's exit
					return
				}
				continue
			}
			// call the user map function
			// TODO: add a recover here to gracefully recover from failed map
			mapRes := mapf(strconv.Itoa(reply.Job.MapJob.JobId), string(content))
			// partition the intermediate responses into files
			// TODO: optimize writing to partition
			mapJobPartitions := map[int]string{}
			for _, kv := range mapRes {
				partition := ihash(kv.Key) % reply.Job.MapJob.nReduce
				ofile, err := os.OpenFile(fmt.Sprintf("map-%d-%d", reply.Job.MapJob.JobId, partition), os.O_CREATE | os.O_WRONLY | os.O_APPEND, 0644)
				if err != nil {
					log.Printf("cannot open file to write partition %d of map task %d\n", partition, reply.Job.MapJob.JobId)
					doneArgs.err = fmt.Errorf("cannot open file to write partition %d of map task %d", partition, reply.Job.MapJob.JobId)
					ok = call("Coordinator.JobDone", &doneArgs, &doneReply)
					if !ok {
						// coordinator is done let's exit
						return
					}
					break
				}
				// write out the intermediate kv as a json string to the file
				jsonEnc := json.NewEncoder(ofile)
				jsonEnc.Encode(kv)
				mapJobPartitions[partition] = ofile.Name()
			}
			if doneArgs.err != nil {
				// if we already reported an error, just ask for a new task
				continue
			}
			// build the map job partitions
			var jobPartitions MapJobIntermediateFiles[int]
			for partition, fileName := range mapJobPartitions {
				jobPartitions.IntermediateFiles = append(jobPartitions.IntermediateFiles, IntermediateFile[int]{PartitionId: partition, FileName: fileName})
			}
			// respond to the coordinator
			doneArgs = JobDoneReq[int]{
				JobType:         Map,
				Job: reply.Job,
				MapJobPartitions: jobPartitions,
			}
		} else {
			// read all intermediate files
			kva := []KeyValue{}
			for _, file := range reply.Files {
				// make sure that if we had recorded an error while reading files
				// break the loop
				if doneArgs.err != nil {
					break
				}
				fileD, err := os.Open(file)
				if err != nil {
					log.Printf("reduce task %d cannot open intermediate file %s\n", reply.Job.ReduceJob.JobId, file)
					doneArgs.err = fmt.Errorf("reduce task %d cannot open intermediate file %s", reply.Job.ReduceJob.JobId, file)
					ok = call("Coordinator.JobDone", &doneArgs, &doneReply)
					if !ok {
						return
					}
					break
				}

				jsonDec := json.NewDecoder(fileD)
				for {
					var kv KeyValue
					if err := jsonDec.Decode(&kv); err != nil {
						log.Printf("cannot read json encoded intermediate file %s", file)
						doneArgs.err = fmt.Errorf("cannot read json encoded intermediate file %s", file)
						ok = call("Coordinator.JobDone", &doneArgs, &doneReply)
						if !ok {
							return
						}
						break
					}
					kva = append(kva, kv)
				}
			}
			if doneArgs.err != nil {
				continue
			}
			// sort the intermediate keys because multiple keys can map to
			// a partition (reducer job id) but we want to send unique
			// keys and collated values for those keys to the reduce function
			// 	map (k1,v1) → list(k2,v2)
			// reduce (k2,list(v2)) → list(v2)
			sort.Sort(ByKey(kva))

			oname := fmt.Sprintf("mr-out-%s", reply.Job.ReduceJob.JobId)
			ofile, _ := os.Create(oname)

			//
			// call Reduce on each distinct key in intermediate[],
			// and print the result to mr-out-0.
			//
			i := 0
			for i < len(kva) {
				j := i + 1
				for j < len(kva) && kva[j].Key == kva[i].Key {
					j++
				}
				values := []string{}
				for k := i; k < j; k++ {
					values = append(values, kva[k].Value)
				}
				output := reducef(kva[i].Key, values)

				// this is the correct format for each line of Reduce output.
				fmt.Fprintf(ofile, "%v %v\n", kva[i].Key, output)

				i = j
			}
			// reply to the coordinator
			doneArgs = JobDoneReq[int]{
				JobType: Reduce,
				Job: reply.Job,
				MapJobPartitions: MapJobIntermediateFiles[int]{},
				err:   nil,
			}
		}
		ok = call("Coordinator.JobDone", &doneArgs, &doneReply)
		if !ok {
			return
		}
	}
}

// example function to show how to make an RPC call to the coordinator.
//
// the RPC argument and reply types are defined in rpc.go.
func CallExample() {

	// declare an argument structure.
	args := ExampleArgs{}

	// fill in the argument(s).
	args.X = 99

	// declare a reply structure.
	reply := ExampleReply{}

	// send the RPC request, wait for the reply.
	// the "Coordinator.Example" tells the
	// receiving server that we'd like to call
	// the Example() method of struct Coordinator.
	ok := call("Coordinator.Example", &args, &reply)
	if ok {
		// reply.Y should be 100.
		fmt.Printf("reply.Y %v\n", reply.Y)
	} else {
		fmt.Printf("call failed!\n")
	}
}

// send an RPC request to the coordinator, wait for the response.
// usually returns true.
// returns false if something goes wrong.
func call(rpcname string, args interface{}, reply interface{}) bool {
	// c, err := rpc.DialHTTP("tcp", "127.0.0.1"+":1234")
	sockname := coordinatorSock()
	c, err := rpc.DialHTTP("unix", sockname)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	defer c.Close()

	err = c.Call(rpcname, args, reply)
	if err == nil {
		return true
	}

	fmt.Println(err)
	return false
}
