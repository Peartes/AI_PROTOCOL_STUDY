// this main file starts a kv server that uses Raft
// each kv server communicates with a specific raft server
package main

import (
	"log"
	"net/rpc"
	"os"
	"strconv"

	"6.5840/kvsrv"
	"6.5840/labrpc"
	"6.5840/raft"
)

var raftServerId int64
var err error

func init() {
	if len(os.Args) > 1 {
		raftServerId, err = strconv.ParseInt(os.Args[1], 0, 0)
		if err != nil {
			log.Fatal("cannot not convert raft id into a int64")
		}
	} else {
		log.Fatal("provide the raft instance id as an argument")
	}
}

type RaftEndpoint struct {
	endName string // this end-point's name
}

func main() {
	// make raft servers
	servers := 2
	peers := make([]labrpc.ServiceEndpoint, servers)
	persister := raft.MakePersister()
	for i := 0; i < servers; i++ {
		peers[i] = &RaftEndpoint{endName: raft.CoordinatorSock(i)}
	}
	raft := raft.Make(peers, int(raftServerId), persister, nil)
	// make kv server
	kvsrv := kvsrv.StartKVServer(strconv.FormatInt(raftServerId, 10))
	kvsrv.Raft = raft
	select {} // keep the server running
}

func (re *RaftEndpoint) Call(svcMeth string, args interface{}, reply interface{}) bool {
	c, err := rpc.DialHTTP("unix", re.endName)
	if err != nil {
		log.Printf("dialing:%s failed", re.endName)
		return false
	}
	defer c.Close()

	err = c.Call(svcMeth, args, reply)
	if err == nil {
		return true
	}

	log.Printf("calling:%s failed", re.endName)
	return false
}
