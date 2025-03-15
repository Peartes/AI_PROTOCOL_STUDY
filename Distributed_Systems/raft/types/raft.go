package types

type RequestVoteRequest struct {
	Term         int32
	CandidateId  int32
	LastLogIndex int32
	LastLogTerm  int32
}

type RequestVoteResponse struct {
	Term        int32
	VoteGranted bool
}
