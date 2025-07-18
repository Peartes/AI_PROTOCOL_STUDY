package kvsrv

// Put or Append
type PutAppendArgs struct {
	Key   string
	Value string
	ClientId string
	RequestNumber int
}

type PutAppendReply struct {
	Value string
}

type GetArgs struct {
	Key string
	ClientId string
	RequestNumber int
}

type GetReply struct {
	Value string
}
