package question

import (
	"context"
	"fmt"
	"quiz-mod/internal/config"
	"quiz-mod/internal/model"
)

// сервис, управляющий вопросами
type Question struct {
	firstLevel  []model.SimpleQuestion
	secondLevel []*model.HardQuestion
	thirdLevel  []model.SimpleQuestion

	users map[int64]userState // для хранения состояний пользователей

	storage storage
}

type storage interface {
	SaveResults(rctx context.Context, es model.Result) error
}

func New(cfg *config.Config, storage storage) *Question {
	return &Question{
		firstLevel:  cfg.FirstLevel,
		secondLevel: cfg.SecondLevel,
		thirdLevel:  cfg.ThirdLevel,
		users:       make(map[int64]userState),
		storage:     storage,
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
