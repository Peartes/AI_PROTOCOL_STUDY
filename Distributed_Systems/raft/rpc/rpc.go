package rpc

import (
	"context"

	"github.com/peartes/distr_system/raft/types"
)

type RaftServer struct {
	types.UnimplementedRaftServerServer
}

var _ types.RaftServerServer = &RaftServer{}

func (r *RaftServer) RequestVote(ctx context.Context, payload *types.RequestVoteRequest) (*types.RequestVoteResponse, error) {
	return &types.RequestVoteResponse{}, nil
}
