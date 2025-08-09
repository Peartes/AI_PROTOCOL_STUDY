package main

import (
	"log"
	"net/rpc"
	"os"
	"reflect"
	"strconv"

	"6.5840/labrpc"
	"6.5840/raft"
)
type reqMsg struct {
	endname  interface{} // name of sending ClientEnd
	svcMeth  string      // e.g. "Raft.AppendEntries"
	argsType reflect.Type
	args     []byte
	replyCh  chan replyMsg
}

type replyMsg struct {
	ok    bool
	reply []byte
}

type RaftEndpoint struct {
	endName string   // this end-point's name
	ch      chan reqMsg   // copy of Network.endCh
	done    chan struct{} // closed when Network is cleaned up
}

var _ labrpc.ServiceEndpoint = &RaftEndpoint{}
var raftId int64
var err error

func init() {
	if len(os.Args) > 1 {
		raftId, err = strconv.ParseInt(os.Args[1], 0, 0)
		if err != nil {
			log.Fatal("cannot not convert raft id into a int64")
		}
	} else {
		log.Fatal("provide the raft instance id as an argument")
	}
}
func main() {
	servers := 3
	peers := make([]labrpc.ServiceEndpoint, servers)
	persister := raft.MakePersister()
	for i := 0; i < servers; i++ {
		peers[i] = &RaftEndpoint{endName: raft.CoordinatorSock(i)}
	}
	raft.Make(peers, int(raftId), persister, nil)
	select {}
}

func (re *RaftEndpoint) Call(svcMeth string, args interface{}, reply interface{}) bool {
	c, err := rpc.DialHTTP("unix", re.endName)
	if err != nil {
		log.Print("dialing:%s failed", re.endName)
		return false
	}
	defer c.Close()

	err = c.Call(svcMeth, args, reply)
	if err == nil {
		return true
	}

	log.Print("calling:%s failed", re.endName)
	return false
}
