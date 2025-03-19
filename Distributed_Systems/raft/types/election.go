package types

func (s *State) BecomeLeader() bool {
	s.SetState(Candidate)
	s.SetVotedFor(s.GetServerId())

	return true
}
