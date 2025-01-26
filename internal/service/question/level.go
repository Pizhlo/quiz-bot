package question

import (
	"quiz-bot/internal/model"
	"time"

	"github.com/sirupsen/logrus"
)

func (s *Question) CurrentLevel(userID int64) (int, error) {
	state, err := s.stateByUser(userID)
	if err != nil {
		return 0, err
	}

	return state.level, nil
}

func (s *Question) StartFirstLvl(userID int64) {
	logrus.Debugf("Start first level with user %d", userID)

	state := userState{
		level:        model.FirstLevel,
		rigthAnswers: 0,
		maxQuestions: len(s.firstLevel),
		question:     0,
		startTime:    time.Now(),
	}

	state.result.InitRigthAnswers(userID)
	state.result.InitTotalAnswers(userID)

	state.result.SaveTotalAnswers(userID, model.FirstLevel, len(s.firstLevel))

	s.saveState(userID, state)
}

func (s *Question) StartSecondLvl(userID int64) error {
	logrus.Debugf("Start second level with user %d", userID)

	state, err := s.stateByUser(userID)
	if err != nil {
		return err
	}

	// сохраняем общее количество вопросов на уровне
	state.result.SaveTotalAnswers(userID, model.SecondLevel, len(s.secondLevel))

	state.level = model.SecondLevel
	state.maxQuestions = len(s.secondLevel)
	state.question = 0
	state.rigthAnswers = 0

	s.saveState(userID, state)

	return nil
}

func (s *Question) StartThirdLvl(userID int64) error {
	logrus.Debugf("Start third level with user %d", userID)

	state, err := s.stateByUser(userID)
	if err != nil {
		return err
	}

	state.result.SaveTotalAnswers(userID, model.ThirdLevel, len(s.thirdLevel))

	state.level = model.ThirdLevel
	state.maxQuestions = len(s.thirdLevel)
	state.question = 0
	state.rigthAnswers = 0

	s.saveState(userID, state)

	return nil

}

func (s *Question) LevelResults(userID int64) (int, error) {
	state, err := s.stateByUser(userID)
	if err != nil {
		return 0, err
	}

	return state.rigthAnswers, nil
}
