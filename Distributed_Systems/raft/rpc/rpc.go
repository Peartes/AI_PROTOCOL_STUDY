package rpc

type RequestVoteMessage struct {
	Term         int
	CandidateId  int
	LastLogIndex int
	LastLogTerm  int
}

func RequestVote(payload RequestVoteMessage) bool {
	return false
}
