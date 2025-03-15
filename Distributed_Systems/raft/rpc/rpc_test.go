package rpc_test

import (
	"testing"

	"github.com/peartes/distr_system/raft/rpc"
	"github.com/peartes/distr_system/raft/types"
	"github.com/stretchr/testify/require"
)

func TestRequestVote(t *testing.T) {
	scenarios := []struct {
		Request     *types.RequestVoteRequest
		Response    *types.RequestVoteResponse
		serverState func() *rpc.RaftServer
	}{
		{
			// initial request with all servers just starting up
			Request: &types.RequestVoteRequest{
				Term:         1,
				CandidateId:  2,
				LastLogIndex: 0,
				LastLogTerm:  1,
			},
			Response: &types.RequestVoteResponse{
				Term:        1,
				VoteGranted: true,
			},
			serverState: func() *rpc.RaftServer { return &rpc.RaftServer{State: types.NewState(1)} },
		},
		{
			// stale candidate term
			Request: &types.RequestVoteRequest{
				Term:         1,
				CandidateId:  2,
				LastLogIndex: 1,
				LastLogTerm:  1,
			},
			Response: &types.RequestVoteResponse{
				Term:        3,
				VoteGranted: false,
			},
			serverState: func() *rpc.RaftServer {
				state := types.NewState(1)
				state.SetCurrentTerm(3)

				return &rpc.RaftServer{State: state}
			},
		},
		{
			// stale candidate log index
			Request: &types.RequestVoteRequest{
				Term:         1,
				CandidateId:  2,
				LastLogIndex: 0,
				LastLogTerm:  1,
			},
			Response: &types.RequestVoteResponse{
				Term:        1,
				VoteGranted: false,
			},
			serverState: func() *rpc.RaftServer {
				state := types.NewState(1)
				state.SetLog([]types.Log{{Command: "", Term: 1}, {Command: "", Term: 1}})

				return &rpc.RaftServer{State: state}
			},
		},
		{
			// stale candidate log term
			Request: &types.RequestVoteRequest{
				Term:         1,
				CandidateId:  2,
				LastLogIndex: 1,
				LastLogTerm:  1,
			},
			Response: &types.RequestVoteResponse{
				Term:        2,
				VoteGranted: false,
			},
			serverState: func() *rpc.RaftServer {
				state := types.NewState(1)
				state.SetCurrentTerm(2)
				state.SetLog([]types.Log{{Command: "", Term: 1}, {Command: "", Term: 2}})

				return &rpc.RaftServer{State: state}
			},
		},
		{
			// server already voted
			Request: &types.RequestVoteRequest{
				Term:         1,
				CandidateId:  2,
				LastLogIndex: 1,
				LastLogTerm:  1,
			},
			Response: &types.RequestVoteResponse{
				Term:        1,
				VoteGranted: false,
			},
			serverState: func() *rpc.RaftServer {
				state := types.NewState(1)
				state.SetVotedFor(3)

				return &rpc.RaftServer{State: state}
			},
		},
	}

	for _, tt := range scenarios {
		t.Run("should request vote", func(t *testing.T) {
			res := &types.RequestVoteResponse{}
			tt.serverState().RequestVote(tt.Request, res)
			require.EqualValues(t, tt.Response, res)
		})
	}
}
