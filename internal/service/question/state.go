package question

import "fmt"

const (
	firstLevel = iota
	secondLevel
	thirdLevel
)

type userState struct {
	level        int
	question     int
	maxQuestions int
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
