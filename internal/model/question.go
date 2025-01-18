package model

import (
	"fmt"
	"quiz-mod/internal/message"
	"reflect"

	"github.com/sirupsen/logrus"
)

type Question struct {
	Text    string   `json:"question"`
	Answers []string `json:"answers"`
}

func (q *Question) QuestionText(currIdx, maxIdx int) string {
	msg := message.Question

	return fmt.Sprintf(msg, currIdx, maxIdx, q.Text)
}

// вопрос, у которого только один правильный ответ
type SimpleQuestion struct {
	Question
	RigthAnswer string `json:"rigth_answer"`
	UserAnswer  map[int64]string
}

func (s *SimpleQuestion) SetUserAnswer(user int64, answer string) {
	if s.UserAnswer == nil {
		s.UserAnswer = make(map[int64]string)
	}

	logrus.Debugf("Set user's answer. User: %+v. Answer: %+v", user, answer)

	s.UserAnswer[user] = answer
}

func (s *SimpleQuestion) Valid(answer string) bool {
	return s.RigthAnswer == answer
}

// вопрос, у которого несколько правильных ответов
type HardQuestion struct {
	Question
	RigthAnswers []string `json:"rigth_answers"`
	UserAnswers  map[int64][]string
}

func (s *HardQuestion) AddUserAnswer(user int64, answer string) {
	if s.UserAnswers == nil {
		s.UserAnswers = make(map[int64][]string)
	}

	logrus.Debugf("Added user's answer. User: %+v. Answer: %+v", user, answer)

	s.UserAnswers[user] = append(s.UserAnswers[user], answer)
}

func (s *HardQuestion) Valid(userID int64) bool {
	answers := s.UserAnswers[userID]
	return reflect.DeepEqual(answers, s.RigthAnswers)
}
