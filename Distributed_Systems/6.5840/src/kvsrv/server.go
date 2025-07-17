package kvsrv

import (
	"fmt"
	"log"
	"math"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"sync"
)

const Debug = false

func DPrintf(format string, a ...interface{}) (n int, err error) {
	if Debug {
		log.Printf(format, a...)
	}
	return
}


type KVServer struct {
	mu sync.Mutex

	state map[string]string
	// this is the mapping of the requestid (clientid+req_num) -> response in order to not reprocess a retransmitted request
	// which will break consistency (linearizability)
	processed map[string]string
}


func (kv *KVServer) Get(args *GetArgs, reply *GetReply) error {
	// TODO: check if this read is a re-transmission and send back old result
	reply.Value = kv.state[args.Key]
	return nil
}

func (kv *KVServer) Put(args *PutAppendArgs, reply *PutAppendReply) error {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	requestId := fmt.Sprintf("%s%d", args.ClientId, args.RequestNumber)
	oldRequestId := fmt.Sprintf("%s%f", args.ClientId, math.Max(float64(args.RequestNumber - 1), 0))
	value, ok := kv.processed[requestId]
	if !ok {
		// this is a new request so let's process it
		// send back the new value of this key
		reply.Value = args.Value
		kv.state[args.Key] = args.Value
		// save this in our processed state in case of re-transmissions
		delete(kv.processed, oldRequestId)
		kv.processed[requestId] = args.Value
	} else  {
		// we have processed this request, return back the value
		reply.Value = value
	}
	return nil
}

func (kv *KVServer) Append(args *PutAppendArgs, reply *PutAppendReply) error {
	// to guarantee linearizability for put, we need to consider 2 conditions
	// the write hasn't been done previously
	// we can guarantee that by saving previous writes as to not repeat them again because
	// when old writes are re-processed (even though the system is idempotent) it could break consistency
	// the second thing is that we don't forget already completed writes during a crash
	// we can achieve that by saving the completed write to disk and on boot-up, we read state from disk
	kv.mu.Lock()
	defer kv.mu.Unlock()
	requestId := fmt.Sprintf("%s%d", args.ClientId, args.RequestNumber)
	oldRequestId := fmt.Sprintf("%s%f", args.ClientId, math.Max(float64(args.RequestNumber - 1), 0))
	value, ok := kv.processed[requestId]
	if !ok {
		// this is a new request so let's process it
		// send back the new value of this key
		reply.Value = args.Value
		oldValue := kv.state[args.Key]
		kv.state[args.Key] = fmt.Sprintf("%s%s", oldValue, args.Value)
		// save this in our processed state in case of re-transmissions
		delete(kv.processed, oldRequestId)
		kv.processed[requestId] = args.Value
	} else  {
		// we have processed this request, return back the value
		reply.Value = value
	}
	return nil
}

func StartKVServer() *KVServer {
	kv := new(KVServer)

	// You may need initialization code here.
	// create a listener on a Unix domain socket
	// and start the RPC server
	rpc.RegisterName("KVServer", kv)
	rpc.HandleHTTP()
	//l, e := net.Listen("tcp", ":1234")
	sockname := coordinatorSock()
	os.Remove(sockname)
	l, e := net.Listen("unix", sockname)
	if e != nil {
		DPrintf("listen error:", e)
	}
	go http.Serve(l, nil)

	return kv
}
