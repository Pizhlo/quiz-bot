package question

import "github.com/sirupsen/logrus"

func (s *Question) StartFirstLvl(userID int64) {
	logrus.Debugf("Start first level with user %d", userID)

	s.users[userID] = userState{}
	s.users[userID] = userState{
		level:        firstLevel,
		maxQuestions: len(s.firstLevel),
		question:     0,
	}
}
