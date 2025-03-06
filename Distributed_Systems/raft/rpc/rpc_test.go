package rpc_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/peartes/distr_system/raft/rpc"
	"github.com/peartes/distr_system/raft/types"
	"github.com/stretchr/testify/require"
)

func TestRequestVote(t *testing.T) {
	scenarios := []struct {
		Request     *types.RequestVoteRequest
		err         error
		Response    *types.RequestVoteResponse
		serverState func() *types.State
	}{
		{
			// initial request with all servers just starting up
			Request: &types.RequestVoteRequest{
				Term:         1,
				CandidateId:  2,
				LastLogIndex: 0,
				LastLogTerm:  0,
			},
			err: nil,
			Response: &types.RequestVoteResponse{
				Term:        1,
				VoteGranted: true,
			},
			serverState: func() *types.State { return types.NewState(1) },
		},
		{
			// stale candidate term
			Request: &types.RequestVoteRequest{
				Term:         1,
				CandidateId:  2,
				LastLogIndex: 1,
				LastLogTerm:  1,
			},
			err:      fmt.Errorf("candidate term %d is lower than server term %d", 1, 2),
			Response: nil,
			serverState: func() *types.State {
				state := types.NewState(1)
				state.SetCurrentTerm(2)

				return state
			},
		},
		{
			// stale candidate log index
			Request: &types.RequestVoteRequest{
				Term:         1,
				CandidateId:  2,
				LastLogIndex: 1,
				LastLogTerm:  1,
			},
			err:      fmt.Errorf("candidate log term %d and index %d is not as up to date as server log term %d and index %d", 1, 1, 1, 2),
			Response: nil,
			serverState: func() *types.State {
				state := types.NewState(1)
				state.SetLog([]types.Log{{Command: "", Term: 1}, {Command: "", Term: 1}})

				return state
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
			err:      fmt.Errorf("candidate log term %d and index %d is not as up to date as server log term %d and index %d", 1, 1, 1, 2),
			Response: nil,
			serverState: func() *types.State {
				state := types.NewState(1)
				state.SetCurrentTerm(2)
				state.SetLog([]types.Log{{Command: "", Term: 1}, {Command: "", Term: 2}})

				return state
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
			err: nil,
			Response: &types.RequestVoteResponse{
				Term:        1,
				VoteGranted: false,
			},
			serverState: func() *types.State {
				state := types.NewState(1)
				state.SetVotedFor(3)

				return state
			},
		},
	}

	server := rpc.RaftServer{
		types.UnimplementedRaftServerServer{},
	}

	for _, tt := range scenarios {
		t.Run("should request vote", func(t *testing.T) {
			ctx := context.WithValue(context.Background(), types.NodeKey, tt.serverState())

			res, err := server.RequestVote(ctx, tt.Request)

			if tt.err != nil {
				require.Error(t, tt.err, err)
			} else {
				require.EqualValues(t, &tt.Response, &res)
			}
		})
	}
}
