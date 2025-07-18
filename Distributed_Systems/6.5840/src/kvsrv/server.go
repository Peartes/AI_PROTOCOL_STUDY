package kvsrv

import (
	"fmt"
	"log"
	"sync"
)

const Debug = true

func DPrintf(format string, a ...any) (n int, err error) {
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

func (kv *KVServer) Get(args *GetArgs, reply *GetReply) {

	kv.mu.Lock()
	defer kv.mu.Unlock()
	reply.Value = kv.state[args.Key]
}

func (kv *KVServer) Put(args *PutAppendArgs, reply *PutAppendReply) {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	requestId := fmt.Sprintf("%s%d", args.ClientId, args.RequestNumber)
	oldRequestId := fmt.Sprintf("%s%d", args.ClientId, args.RequestNumber-1)
	value, ok := kv.processed[requestId]
	if !ok {
		// this is a new request so let's process it
		// oldValue := kv.state[args.Key]
		// send back the new value of this key
		reply.Value = args.Value
		kv.state[args.Key] = args.Value
		// save this in our processed state in case of re-transmissions
		delete(kv.processed, oldRequestId)
		kv.processed[requestId] = args.Value
	} else {
		// we have processed this request, return back the value
		reply.Value = value
	}
}

func (kv *KVServer) Append(args *PutAppendArgs, reply *PutAppendReply) {
	// to guarantee linearizability for put, we need to consider 2 conditions
	// the write hasn't been done previously
	// we can guarantee that by saving previous writes as to not repeat them again because
	// when old writes are re-processed (even though the system is idempotent) it could break consistency
	// the second thing is that we don't forget already completed writes during a crash
	// we can achieve that by saving the completed write to disk and on boot-up, we read state from disk
	kv.mu.Lock()
	defer kv.mu.Unlock()
	requestId := fmt.Sprintf("%s%d", args.ClientId, args.RequestNumber)
	oldRequestId := fmt.Sprintf("%s%d", args.ClientId, args.RequestNumber-1)
	value, ok := kv.processed[requestId]
	if !ok {
		// this is a new request so let's process it
		oldValue := kv.state[args.Key]
		// send back the old value of this key
		reply.Value = oldValue
		kv.state[args.Key] = fmt.Sprintf("%s%s", oldValue, args.Value)
		// save this in our processed state in case of re-transmissions
		delete(kv.processed, oldRequestId)
		kv.processed[requestId] = oldValue
	} else {
		// we have processed this request, return back the value
		reply.Value = value
	}
}

func StartKVServer() *KVServer {
	kv := new(KVServer)

	// You may need initialization code here.
	kv.state = map[string]string{}
	kv.processed = map[string]string{}

	return kv
}
