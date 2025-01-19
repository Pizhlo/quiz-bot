package question

import (
	"context"
	"fmt"
	"quiz-mod/internal/model"
	"time"

	"github.com/sirupsen/logrus"
)

func (s *Question) SetAnswer(userID int64, answer string) error {
	state, err := s.stateByUser(userID)
	if err != nil {
		return err
	}

	switch state.level {
	case model.FirstLevel:
		question := s.firstLevel[state.question]

		if question.Valid(answer) {
			state.rigthAnswers++
			logrus.Debugf("set answer first lvl")
			state.result.SaveAnswers(userID, model.FirstLevel, state.rigthAnswers)
			s.saveState(userID, state)
		}
	case model.ThirdLevel:
		question := s.thirdLevel[state.question]

		if question.Valid(answer) {
			state.rigthAnswers++
			logrus.Debugf("set answer third lvl")
			state.result.SaveAnswers(userID, model.ThirdLevel, state.rigthAnswers)
			s.saveState(userID, state)
		}
	default:
		return fmt.Errorf("invalid level for simple question: %+v", state.level)
	}

	return nil
}

func (s *Question) SaveAnswers(userID int64) error {
	state, err := s.stateByUser(userID)
	if err != nil {
		return err
	}

	switch state.level {
	case model.SecondLevel:
		question := s.secondLevel[state.question]

		if question.Valid(userID) {
			state.rigthAnswers++

			state.result.SaveAnswers(userID, model.SecondLevel, state.rigthAnswers)
			s.saveState(userID, state)
		}
	default:
		return fmt.Errorf("invalid level for hard question: %+v", state.level)
	}

	return nil
}

func (s *Question) AddAnswer(userID int64, answer string) error {
	state, err := s.stateByUser(userID)
	if err != nil {
		return err
	}

	switch state.level {
	case model.SecondLevel:
		question := s.secondLevel[state.question]

		question.AddUserAnswer(userID, answer)
	default:
		return fmt.Errorf("invalid level for hard question: %+v", state.level)
	}

	return nil
}

func (s *Question) RigthAnswer(userID int64) ([]string, error) {
	state, err := s.stateByUser(userID)
	if err != nil {
		return nil, err
	}

	switch state.level {
	case model.FirstLevel:
		question := s.firstLevel[state.question]
		return []string{question.RigthAnswer}, nil
	case model.SecondLevel:
		question := s.secondLevel[state.question]
		return question.RigthAnswers, nil
	case model.ThirdLevel:
		question := s.thirdLevel[state.question]
		return []string{question.RigthAnswer}, nil
	default:
		return nil, fmt.Errorf("unknown level: %+v", state.level)
	}
}

func (s *Question) UserAnswers(userID int64) ([]string, error) {
	state, err := s.stateByUser(userID)
	if err != nil {
		return nil, err
	}

	switch state.level {
	case model.SecondLevel:
		question := s.secondLevel[state.question]
		return question.UserAnswers[userID], nil
	default:
		return nil, fmt.Errorf("invalid level for hard question: %+v", state.level)
	}
}

func (s *Question) Results(userID int64) (model.Result, error) {
	state, err := s.stateByUser(userID)
	if err != nil {
		return model.Result{}, err
	}

	res := state.result

	return res, nil
}

// StopTimer записывает, сколько длилась викторина
func (s *Question) StopTimer(userID int64) error {
	state, err := s.stateByUser(userID)
	if err != nil {
		return err
	}

	res := state.result

	res.Duration = time.Since(state.startTime)

	res.Date = time.Now()

	state.result = res

	s.saveState(userID, state)

	return nil
}

// SaveResults сохраняет результаты в БД
func (s *Question) SaveResults(ctx context.Context, userID int64) error {
	res, err := s.Results(userID)
	if err != nil {
		return err
	}

	res.TgID = userID

	// проверяем результаты на валидность перед сохранением
	err = res.Valid()
	if err != nil {
		return err
	}
	return s.storage.SaveResults(ctx, res)
}
