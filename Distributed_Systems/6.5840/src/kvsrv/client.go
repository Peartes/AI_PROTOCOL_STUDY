package kvsrv

import (
	"crypto/rand"
	"math/big"
	"net/rpc"

	"6.5840/labrpc"
)


type Clerk struct {
	server *labrpc.ClientEnd
	requestNumber int // monotonically increasing on every request
	clientId string // unique id for this client
	
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
	args := GetArgs{Key: key}
	reply := &GetReply{}
	err := ck.call("KVServer", &args, reply)
	if !err {
		ck.requestNumber = ck.requestNumber + 1
	}
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
	args := PutAppendArgs{Key: key, Value: value}
	reply := &GetReply{}
	err := ck.call("KVServer."+op, &args, reply)
	if !err {
		ck.requestNumber = ck.requestNumber + 1
	}
	return reply.Value
}

func (ck *Clerk) Put(key string, value string) {

	ck.PutAppend(key, value, "Put")
}

// Append value to key's value and return that value
func (ck *Clerk) Append(key string, value string) string {
	return ck.PutAppend(key, value, "Append")
}

// send an RPC request to the coordinator, wait for the response.
// usually returns true.
// returns false if something goes wrong.
func (ck *Clerk) call(rpcname string, args any, reply any) bool {
	c, err := rpc.DialHTTP("unix", coordinatorSock())
	if err != nil {
		DPrintf("error dialing server %v:", err.Error())
	}
	defer c.Close()

	err = c.Call(rpcname, args, reply)
	if err == nil {
		return true
	}

	DPrintf(err.Error())
	return false
}