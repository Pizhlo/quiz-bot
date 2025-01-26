package question

import (
	"context"
	"fmt"
	"quiz-bot/internal/model"
	"quiz-bot/internal/view"

	"gopkg.in/telebot.v3"
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

func (s *Question) AllResults(ctx context.Context, userID int64) (string, *telebot.ReplyMarkup, error) {
	results, err := s.storage.AllResults(ctx, userID)
	if err != nil {
		return "", nil, err
	}

	view := view.NewResult()

	s.views[userID] = view

	return view.Message(results), view.Keyboard(), nil
}
