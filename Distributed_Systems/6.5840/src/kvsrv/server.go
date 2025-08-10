package kvsrv

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"sync"

	"6.5840/raft"
)

const Debug = true

const ErrNotLeader = "ErrNotLeader"

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
	processed map[string]map[string]string
	applyChan chan raft.ApplyMsg // this is used to notify the kv server that a new log entry has been applied
	Raft      *raft.Raft         // this is the raft instance that this kv server is using
}

func (kv *KVServer) Get(args *GetArgs, reply *GetReply) error {

	kv.mu.Lock()
	defer kv.mu.Unlock()
	reply.Value = kv.state[args.Key]
	return nil
}

func (kv *KVServer) Put(args *PutAppendArgs, reply *PutAppendReply) error {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	requestId := fmt.Sprintf("%s%d", args.ClientId, args.RequestNumber)
	_, ok := kv.processed[args.ClientId][requestId]
	if !ok {
		// this is a new request so let's process it
		// oldValue := kv.state[args.Key]
		// send back the new value of this key
		reply.Value = ""
		kv.state[args.Key] = args.Value
		// save this in our processed state in case of re-transmissions
		kv.processed[args.ClientId] = map[string]string{requestId: ""}
	} else {
		// we have processed this request, return back the value
		reply.Value = ""
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
	value, ok := kv.processed[args.ClientId][requestId]
	if !ok {
		// this is a new request so let's process it
		oldValue := kv.state[args.Key]
		// send back the old value of this key
		reply.Value = oldValue
		kv.state[args.Key] = fmt.Sprintf("%s%s", oldValue, args.Value)
		// save this in our processed state in case of re-transmissions
		kv.processed[args.ClientId] = map[string]string{requestId: oldValue}
	} else {
		// we have processed this request, return back the value
		reply.Value = value
	}
	return nil
}

func (kv *KVServer) AppendReplica(args *PutAppendArgs, reply *PutAppendReply) error {
	// to guarantee linearizability for put, we need to consider 2 conditions
	// the write hasn't been done previously
	// we can guarantee that by saving previous writes as to not repeat them again because
	// when old writes are re-processed (even though the system is idempotent) it could break consistency
	// the second thing is that we don't forget already completed writes during a crash
	// we can achieve that by saving the completed write to disk and on boot-up, we read state from disk
	kv.mu.Lock()
	defer kv.mu.Unlock()
	requestId := fmt.Sprintf("%s%d", args.ClientId, args.RequestNumber)
	value, ok := kv.processed[args.ClientId][requestId]
	if !ok {
		// send request to raft to append the value
		fmt.Printf("sending append request to raft for key %s with value %s\n\n", args.Key, args.Value)
		_, _, isLeader := kv.Raft.Start(args)
		if !isLeader {
			return fmt.Errorf("i'm not the leader, current leader is %d", kv.Raft.GetLeader())
		}
	} else {
		// we have processed this request, return back the value
		reply.Value = value
	}
	return nil
}

func StartKVServer(name string) *KVServer {
	kv := new(KVServer)

	// You may need initialization code here.
	kv.state = map[string]string{}
	kv.processed = map[string]map[string]string{}
	kv.applyChan = make(chan raft.ApplyMsg, 1)
	server(kv, name)
	// this is used to notify the kv server that a new log entry has been applied
	go func() {
		for {
			applyCh := <-kv.applyChan
			// process the log entry here
			// this is where we would apply the log entry to the kv server state
			// for now, we just print the state
			fmt.Printf("KVServer receive state update from raft server for put")
			if args, ok := applyCh.Command.(*PutAppendArgs); ok {
				kv.Append(args, &PutAppendReply{})
				fmt.Printf("KVServer applied put command for key	: %s", args.Key)
			}
		}
	}()

	return kv
}

// create a grpc server and register the kv server on it
func server(kv *KVServer, name string) {
	rpc.RegisterName("KVServer", kv)
	// rpc.HandleHTTP()
	sockname := fmt.Sprintf("/var/tmp/kv-%s.sock", name)
	os.Remove(sockname)
	listener, err := net.Listen("unix", sockname)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	go func() {
		if err := http.Serve(listener, nil); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
}
