package types

type State struct {
	// persistent state on each server
	CurrentTerm int
	VotedFor    int
	Log         []string

	// volatile state on each server
	CommitIndex int
	LastApplied int

	// volatile state on leaders
	NextIndex  []int
	MatchIndex []int
}
