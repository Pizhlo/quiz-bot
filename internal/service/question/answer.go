package question

import (
	"fmt"
	"quiz-bot/internal/model"
	"time"
)

// SetAnswer сохраняет один ответ на вопрос. Используется для простых вопросов,
// у которых один правильный ответ. Проверяет на правильность ответ пользователя, и,
// если он правильный, увеличивает счетчик state.rigthAnswers
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

			s.SaveState(userID, state)
		}
	case model.ThirdLevel:
		question := s.thirdLevel[state.question]

		if question.Valid(answer) {
			state.rigthAnswers++
			s.SaveState(userID, state)
		}
	default:
		return fmt.Errorf("invalid level for simple question: %+v", state.level)
	}

	return nil
}

// SaveAnswers используется для вопросов со множественным ответом.
// Сравнивает накопленные ответы пользователя с правильными, и, если они совпали,
// увеличивает счетчик state.rigthAnswers
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
			s.SaveState(userID, state)
		}
	default:
		return fmt.Errorf("invalid level for hard question: %+v", state.level)
	}

	return nil
}

// AddAnswer сохраняет список ответов пользователя на вопрос со множественным выбором
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
		return err
	}

	res := state.result

	res.Duration = time.Since(state.startTime)

	res.Seconds = res.Duration.Seconds()

	res.Date = time.Now()

	state.result = res

	s.SaveState(userID, state)

	return nil
}
