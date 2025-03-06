package rpc

import (
	"context"
	"fmt"

	"github.com/peartes/distr_system/raft/types"
)

type key int

const (
	nodeKey key = 1
)

type RaftServer struct {
	types.UnimplementedRaftServerServer
}

var _ types.RaftServerServer = &RaftServer{}

func (r *RaftServer) RequestVote(ctx context.Context, payload *types.RequestVoteRequest) (*types.RequestVoteResponse, error) {
	currState := ctx.Value(nodeKey).(*types.State)

	if currState == nil {
		panic("can not read node state")
	}

	currTerm := int32(currState.GetCurrentTerm())
	logIndex := int32(currState.GetLastLogIndex())
	logTerm := int32(currState.GetLastLogTerm())
	if payload.Term < currTerm {
		return nil, fmt.Errorf("candidate term %d is lower than server term %d", payload.Term, currState.CurrentTerm)
	}

	votedFor := currState.GetVotedFor()
	if votedFor == currState.GetServerId() || votedFor == int(payload.CandidateId) {
		if payload.LastLogIndex >= logIndex && payload.LastLogTerm >= logTerm {
			return &types.RequestVoteResponse{
				Term:        currTerm,
				VoteGranted: true,
			}, nil
		} else {
			return nil, fmt.Errorf("candidate log term %d and index %d is not as up to date as server log term %d and index %d", payload.LastLogTerm, payload.LastLogIndex, logTerm, logIndex)
		}
	}
	return &types.RequestVoteResponse{
		Term:        currTerm,
		VoteGranted: false,
	}, nil
}
