package question

import (
	"context"
	"fmt"
	"quiz-mod/internal/model"
	"quiz-mod/internal/view"
)

func (s *Question) IsQuestionLast(userID int64) (bool, error) {
	state, err := s.stateByUser(userID)
	if err != nil {
		return false, err
	}

	return state.isQuestionLast(), nil
}

func (s *Question) QuestionNum(userID int64) (int, error) {
	state, err := s.stateByUser(userID)
	if err != nil {
		return 0, err
	}

	return state.maxQuestions, nil
}

func (s *Question) CurrentQuestion(userID int64) (*model.Question, error) {
	state, err := s.stateByUser(userID)
	if err != nil {
		return nil, err
	}

	switch state.level {
	case model.FirstLevel:
		return &s.firstLevel[state.question].Question, nil
	case model.SecondLevel:
		return &s.secondLevel[state.question].Question, nil
	case model.ThirdLevel:
		return &s.thirdLevel[state.question].Question, nil
	default:
		return nil, fmt.Errorf("invalid level for simple question: %+v", state.level)
	}
}

func (s *Question) AllResults(ctx context.Context, userID int64) (string, error) {
	results, err := s.storage.AllResults(ctx, userID)
	if err != nil {
		return "", err
	}

	view := view.NewResult()

	return view.Message(results), nil
}
