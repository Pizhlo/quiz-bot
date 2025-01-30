package model

import (
	"fmt"
	"quiz-bot/internal/message"
	"reflect"
	"strings"

	"github.com/sirupsen/logrus"
)

const (
	FirstLevel = iota
	SecondLevel
	ThirdLevel
)

type Question struct {
	Text    string   `json:"question"` // текст вопроса
	Answers []string `json:"answers"`  // варианты ответа
	Picture string   `json:"picture"`  // картинка (опционально)
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

func (s *SimpleQuestion) Reset(userID int64) {
	delete(s.UserAnswer, userID)
}

func (s *SimpleQuestion) Valid(answer string) bool {
	return strings.EqualFold(s.RigthAnswer, answer)
}

// вопрос, у которого несколько правильных ответов
type HardQuestion struct {
	Question
	RigthAnswers []string           `json:"rigth_answers"` // правильные ответы
	UserAnswers  map[int64][]string // ответы пользователя (по userID)
}

func (s *HardQuestion) AddUserAnswer(user int64, answer string) {
	if s.UserAnswers == nil {
		s.UserAnswers = make(map[int64][]string)
	}

	logrus.Debugf("Added user's answer. User: %+v. Answer: %+v", user, answer)

	s.UserAnswers[user] = append(s.UserAnswers[user], answer)
}

func (s *HardQuestion) Reset(user int64) {
	delete(s.UserAnswers, user)
}

func (s *HardQuestion) Valid(userID int64) bool {
	answers := s.UserAnswers[userID]
	return reflect.DeepEqual(answers, s.RigthAnswers)
}
