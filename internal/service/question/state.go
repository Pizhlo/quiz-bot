package question

import (
	"fmt"
	"quiz-mod/internal/model"
	"time"
)

type userState struct {
	level        int // текущий уровень
	question     int // текущий вопрос
	maxQuestions int // всего вопросов
	rigthAnswers int // правильных ответов
	startTime    time.Time
	result       model.Result
}

// isQuestionLast возвращает true, если текущий вопрос - последний в раунде
func (s *userState) isQuestionLast() bool {
	return s.question == s.maxQuestions-1
}

func (s *Question) stateByUser(userID int64) (userState, error) {
	state, ok := s.users[userID]
	if !ok {
		return userState{}, fmt.Errorf("not found user's state by user ID")
	}

	return state, nil
}

func (s *Question) saveState(userID int64, state userState) {
	s.users[userID] = state
}

func (s *Question) RigthAnswers(userID int64) (int, error) {
	state, err := s.stateByUser(userID)
	if err != nil {
		return 0, err
	}

	return state.rigthAnswers, nil
}
