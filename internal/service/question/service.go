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

func (s *Question) RigthAnswer(userID int64) ([]string, error) {
	state, err := s.stateByUser(userID)
	if err != nil {
		return nil, err
	}

	switch state.level {
	case firstLevel:
		question := s.firstLevel[state.question]
		return []string{question.RigthAnswer}, nil
	case secondLevel:
		question := s.secondLevel[state.question]
		return question.RigthAnswers, nil
	case thirdLevel:
		question := s.thirdLevel[state.question]
		return []string{question.RigthAnswer}, nil
	default:
		return nil, fmt.Errorf("unknown level: %+v", state.level)
	}
}

func (s *Question) SetNext(userID int64) error {
	state, err := s.stateByUser(userID)
	if err != nil {
		return err
	}

	if state.lastQuestion() {
		state.level++
		state.question = 0
	} else {
		state.question++
	}

	s.saveState(userID, state)

	return nil
}
