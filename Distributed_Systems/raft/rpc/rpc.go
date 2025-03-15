package rpc

import (
	"fmt"

	"github.com/peartes/distr_system/raft/types"
)

type RaftServer struct {
	State *types.State
}

func (r *RaftServer) RequestVote(payload *types.RequestVoteRequest, res *types.RequestVoteResponse) {
	currState := r.State

	if currState == nil {
		panic("can not read node state")
	}

	currTerm := int32(currState.GetCurrentTerm())
	logIndex := int32(currState.GetLastLogIndex())
	logTerm := int32(currState.GetLastLogTerm())
	if payload.Term < currTerm {
		fmt.Printf("candidate term %d is lower than server term %d \n", payload.Term, currState.CurrentTerm)
		res.Term = currTerm
		res.VoteGranted = false
		return
	}

	votedFor := currState.GetVotedFor()
	if votedFor == currState.GetServerId() {
		if payload.LastLogIndex >= logIndex && payload.LastLogTerm >= logTerm {
			currState.SetVotedFor(int(payload.CandidateId))
			res.Term = currTerm
			res.VoteGranted = true
		} else {
			fmt.Printf("candidate %d log term %d and index %d is not as up to date as server log term %d and index %d \n", payload.CandidateId, payload.LastLogTerm, payload.LastLogIndex, logTerm, logIndex)
			res.Term = currTerm
			res.VoteGranted = false
		}
	} else {
		fmt.Printf("server already voted for candidate %d \n", votedFor)
		res.Term = currTerm
		res.VoteGranted = false
	}
}
