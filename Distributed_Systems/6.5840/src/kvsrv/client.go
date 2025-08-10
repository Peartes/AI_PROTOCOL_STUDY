package kvsrv

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"net/rpc"
	"strconv"
	"sync"

	"6.5840/labrpc"
)

type Clerk struct {
	server        *labrpc.ClientEnd
	requestNumber int    // monotonically increasing on every request
	clientId      string // unique id for this client
	mu            sync.Mutex
}

func nrand() int64 {
	max := big.NewInt(int64(1) << 62)
	bigx, _ := rand.Int(rand.Reader, max)
	x := bigx.Int64()
	return x
}

func MakeClerk(server *labrpc.ClientEnd) *Clerk {
	ck := new(Clerk)
	ck.server = server
	ck.clientId = strconv.Itoa(int(nrand()))
	ck.requestNumber = 1
	// You'll have to add code here.
	return ck
}

// fetch the current value for a key.
// returns "" if the key does not exist.
// keeps trying forever in the face of all other errors.
//
// you can send an RPC with code like this:
// ok := ck.server.Call("KVServer.Get", &args, &reply)
//
// the types of args and reply (including whether they are pointers)
// must match the declared types of the RPC handler function's
// arguments. and reply must be passed as a pointer.
func (ck *Clerk) Get(key string) string {
	ck.mu.Lock()
	defer ck.mu.Unlock()
	args := GetArgs{Key: key, ClientId: ck.clientId, RequestNumber: ck.requestNumber}
	reply := &GetReply{}
	done := ck.server.Call("KVServer.Get", &args, reply)

	for !done {
		reply = &GetReply{}
		done = ck.server.Call("KVServer.Get", &args, reply)
	}
	ck.requestNumber = ck.requestNumber + 1
	return reply.Value
}

func (ck *Clerk) GetReplica(key string, server string) string {
	ck.mu.Lock()
	defer ck.mu.Unlock()
	args := GetArgs{Key: key, ClientId: ck.clientId, RequestNumber: ck.requestNumber}
	reply := &GetReply{}
	done := ck.call("KVServer.GetReplica", &args, reply, server)

	for !done {
		reply = &GetReply{}
		done = ck.call("KVServer.GetReplica", &args, reply, server)
	}
	ck.requestNumber = ck.requestNumber + 1
	return reply.Value
}

// shared by Put and Append.
//
// you can send an RPC with code like this:
// ok := ck.server.Call("KVServer."+op, &args, &reply)
//
// the types of args and reply (including whether they are pointers)
// must match the declared types of the RPC handler function's
// arguments. and reply must be passed as a pointer.
func (ck *Clerk) PutAppend(key string, value string, op string) string {
	ck.mu.Lock()
	defer ck.mu.Unlock()
	args := PutAppendArgs{Key: key, Value: value, ClientId: ck.clientId, RequestNumber: ck.requestNumber}
	reply := &GetReply{}
	done := ck.server.Call("KVServer."+op, &args, reply)

	for !done {
		reply = &GetReply{}
		done = ck.server.Call("KVServer."+op, &args, reply)

	}
	ck.requestNumber = ck.requestNumber + 1
	return reply.Value
}

func (ck *Clerk) Put(key string, value string) {

	ck.PutAppend(key, value, "Put")
}

// Append value to key's value and return that value
func (ck *Clerk) Append(key string, value string) string {
	return ck.PutAppend(key, value, "Append")
}

func (ck *Clerk) AppendReplica(key string, value string, server string) string {
	ck.mu.Lock()
	defer ck.mu.Unlock()
	args := PutAppendArgs{Key: key, Value: value, ClientId: ck.clientId, RequestNumber: ck.requestNumber}
	reply := &GetReply{}
	done := ck.call("KVServer.AppendReplica", &args, reply, server)

	for !done {
		reply = &GetReply{}
		done = ck.call("KVServer.AppendReplica", &args, reply, server)
	}
	ck.requestNumber = ck.requestNumber + 1
	return reply.Value
}

func (ck *Clerk) call(svcMeth string, args interface{}, reply interface{}, server string) bool {
	sockName := fmt.Sprintf("/var/tmp/kv-%s.sock", server)
	c, err := rpc.DialHTTP("unix", sockName)
	if err != nil {
		log.Printf("dialing:%s failed", sockName)
		return false
	}
	defer c.Close()

	err = c.Call(svcMeth, args, reply)
	if err == nil {
		return true
	}

	log.Printf("calling:%s failed", server)
	return false
}