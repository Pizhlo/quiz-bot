package question

import (
	"fmt"
	"quiz-mod/internal/config"
	"quiz-mod/internal/model"
)

// сервис, управляющий вопросами
type Question struct {
	firstLevel  []model.SimpleQuestion
	secondLevel []model.HardQuestion
	thirdLevel  []model.SimpleQuestion

	users map[int64]userState // для хранения состояний пользователей
}

const (
	firstLevel = iota
	secondLevel
	thirdLevel
)

type userState struct {
	level    int
	question int
}

func New(cfg *config.Config) *Question {
	return &Question{
		firstLevel:  cfg.FirstLevel,
		secondLevel: cfg.SecondLevel,
		thirdLevel:  cfg.ThirdLevel,
		users:       make(map[int64]userState),
	}
}

func (s *Question) CurrentQuestion(userID int64) (*model.Question, error) {
	state, err := s.stateByUser(userID)
	if err != nil {
		return nil, err
	}

	switch state.level {
	case firstLevel:
		return &s.firstLevel[state.question].Question, nil
	case secondLevel:
		return &s.secondLevel[state.question].Question, nil
	case thirdLevel:
		return &s.thirdLevel[state.question].Question, nil
	default:
		return nil, fmt.Errorf("invalid level for simple question: %+v", state.level)
	}
}

func (s *Question) stateByUser(userID int64) (userState, error) {
	state, ok := s.users[userID]
	if !ok {
		return userState{}, fmt.Errorf("not found user's state by user ID")
	}

	return state, nil
}

func (s *Question) Message(userID int64) (string, error) {
	state, err := s.stateByUser(userID)
	if err != nil {
		return "", err
	}

	curIdx := state.question + 1

	switch state.level {
	case firstLevel:
		maxIdx := len(s.firstLevel)

		question := s.firstLevel[state.question]
		return question.QuestionText(curIdx, maxIdx), nil
	case secondLevel:
		maxIdx := len(s.secondLevel)

		question := s.secondLevel[state.question]
		return question.QuestionText(curIdx, maxIdx), nil
	case thirdLevel:
		maxIdx := len(s.thirdLevel)

		question := s.thirdLevel[state.question]
		return question.QuestionText(curIdx, maxIdx), nil
	default:
		return "", fmt.Errorf("unknown level: %+v", state.level)
	}
}
