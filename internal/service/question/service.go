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

func (s *Question) CurrentLevel(userID int64) (int, error) {
	state, err := s.stateByUser(userID)
	if err != nil {
		return 0, err
	}

	return state.level, nil
}

func (s *Question) Message(userID int64) (string, error) {
	state, err := s.stateByUser(userID)
	if err != nil {
		return "", err
	}

	curIdx := state.question + 1

	switch state.level {
	case firstLevel:
		question := s.firstLevel[state.question]
		return question.QuestionText(curIdx, state.maxQuestions), nil
	case secondLevel:
		question := s.secondLevel[state.question]
		return question.QuestionText(curIdx, state.maxQuestions), nil
	case thirdLevel:
		question := s.thirdLevel[state.question]
		return question.QuestionText(curIdx, state.maxQuestions), nil
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

	if state.isQuestionLast() {
		state.level++
		state.question = 0
	} else {
		state.question++
	}

	s.saveState(userID, state)

	return nil
}

func (s *Question) IsQuestionLast(userID int64) (bool, error) {
	state, err := s.stateByUser(userID)
	if err != nil {
		return false, err
	}

	return state.isQuestionLast(), nil
}
