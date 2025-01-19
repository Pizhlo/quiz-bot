package question

import (
	"context"
	"fmt"
	"quiz-mod/internal/model"
	"time"
)

func (s *Question) SetAnswer(userID int64, answer string) error {
	state, err := s.stateByUser(userID)
	if err != nil {
		return err
	}

	switch state.level {
	case firstLevel:
		question := s.firstLevel[state.question]

		if question.Valid(answer) {
			state.rigthAnswers++
			state.result.SaveAnswers(firstLevel, state.rigthAnswers)
			s.saveState(userID, state)
		}
	case thirdLevel:
		question := s.thirdLevel[state.question]

		if question.Valid(answer) {
			state.rigthAnswers++
			state.result.SaveAnswers(thirdLevel, state.rigthAnswers)
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
	case secondLevel:
		question := s.secondLevel[state.question]

		if question.Valid(userID) {
			state.rigthAnswers++

			state.result.SaveAnswers(secondLevel, state.rigthAnswers)
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
	case secondLevel:
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

func (s *Question) UserAnswers(userID int64) ([]string, error) {
	state, err := s.stateByUser(userID)
	if err != nil {
		return nil, err
	}

	switch state.level {
	case secondLevel:
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

	return s.storage.SaveResults(ctx, res)
}
