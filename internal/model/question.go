package model

import (
	"fmt"
	"quiz-mod/internal/message"
	"reflect"
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
	s.UserAnswer[user] = answer
}

func (s *SimpleQuestion) Valid(user int64) (bool, error) {
	answer, ok := s.UserAnswer[user]
	if !ok {
		return false, fmt.Errorf("not found user's answer")
	}

	return s.RigthAnswer == answer, nil
}

// вопрос, у которого несколько правильных ответов
type HardQuestion struct {
	Question
	RigthAnswers []string `json:"rigth_answers"`
	UserAnswers  map[int64][]string
}

func (s *HardQuestion) SetUserAnswer(user int64, answer []string) {
	s.UserAnswers[user] = answer
}

func (s *HardQuestion) Valid(user int64) (bool, error) {
	answer, ok := s.UserAnswers[user]
	if !ok {
		return false, fmt.Errorf("not found user's answer")
	}

	return reflect.DeepEqual(answer, s.RigthAnswers), nil
}
