package question

import (
	"context"
	"quiz-bot/internal/model"
)

func (s *Question) Results(userID int64) (model.Result, error) {
	state, err := s.stateByUser(userID)
	if err != nil {
		return model.Result{}, err
	}

	res := state.result

	return res, nil
}

// SaveLvlResults сохраняет количество правильных ответов за уровень в результат
func (s *Question) SaveLvlResults(userID int64) error {
	state, err := s.stateByUser(userID)
	if err != nil {
		return err
	}

	state.result.SaveAnswers(state.level, state.rigthAnswers)

	s.SaveState(userID, state)

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
