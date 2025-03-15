package types

func (s *State) BecomeLeader() bool {
	s.SetState(Candidate)
	s.SetVotedFor(s.GetServerId())
	s.SetCurrentTerm(s.GetCurrentTerm() + 1)

	return true
}
