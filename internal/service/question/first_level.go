package question

func (s *Question) StartFirstLvl(userID int64) {
	s.users[userID] = userState{}
	s.users[userID] = userState{
		level:    firstLevel,
		question: 0,
	}
}
