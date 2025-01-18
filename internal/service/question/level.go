package question

import "github.com/sirupsen/logrus"

func (s *Question) StartFirstLvl(userID int64) {
	logrus.Debugf("Start first level with user %d", userID)

	s.users[userID] = userState{
		level:        firstLevel,
		maxQuestions: len(s.firstLevel),
		question:     0,
	}
}

func (s *Question) StartSecondLvl(userID int64) {
	logrus.Debugf("Start second level with user %d", userID)

	s.users[userID] = userState{
		level:        secondLevel,
		maxQuestions: len(s.secondLevel),
		question:     0,
	}
}

func (s *Question) StartThirdLvl(userID int64) {
	logrus.Debugf("Start third level with user %d", userID)

	s.users[userID] = userState{
		level:        thirdLevel,
		maxQuestions: len(s.thirdLevel),
		question:     0,
	}
}

func (s *Question) LevelResults(userID int64) (int, error) {
	state, err := s.stateByUser(userID)
	if err != nil {
		return 0, err
	}

	return state.rigthAnswers, nil
}
