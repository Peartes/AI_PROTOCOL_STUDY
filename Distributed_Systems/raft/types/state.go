package types

type State struct {
	// persistent state on each server
	ServerId    int
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

func NewState() *State {
	return &State{
		CurrentTerm: 0,
		VotedFor:    -1,
		Log:         []string{},
		CommitIndex: 0,
		LastApplied: 0,
		NextIndex:   []int{},
		MatchIndex:  []int{},
	}
}

func (s *State) GetServerId() int {
	return s.ServerId
}

func (s *State) GetLastLogIndex() int {
	return len(s.Log) - 1
}

func (s *State) GetLastLogTerm() int {
	return len(s.Log)
}

func (s *State) GetCommitIndex() int {
	return s.CommitIndex
}

func (s *State) GetLastApplied() int {
	return s.LastApplied
}

func (s *State) GetNextIndex() []int {

	return s.NextIndex
}

func (s *State) GetMatchIndex() []int {
	return s.MatchIndex
}

func (s *State) GetCurrentTerm() int {
	return s.CurrentTerm
}

func (s *State) GetVotedFor() int {
	return s.VotedFor
}

func (s *State) SetCurrentTerm(term int) {
	s.CurrentTerm = term
}

func (s *State) SetVotedFor(votedFor int) {
	s.VotedFor = votedFor
}

func (s *State) SetCommitIndex(commitIndex int) {
	s.CommitIndex = commitIndex
}

func (s *State) SetLastApplied(lastApplied int) {
	s.LastApplied = lastApplied
}

func (s *State) SetNextIndex(nextIndex []int) {
	s.NextIndex = nextIndex
}

func (s *State) SetMatchIndex(matchIndex []int) {
	s.MatchIndex = matchIndex
}

func (s *State) SetLog(log []string) {
	s.Log = log
}

func (s *State) AppendLog(log string) {
	s.Log = append(s.Log, log)
}

func (s *State) GetLog() []string {
	return s.Log
}

func (s *State) GetLogEntry(index int) string {
	return s.Log[index]
}

func (s *State) GetLogEntries(startIndex int) []string {
	return s.Log[startIndex:]
}

func (s *State) GetLogLength() int {
	return len(s.Log)
}
