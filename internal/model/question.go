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
	UserAnswer  string
}

func (s *SimpleQuestion) SetUserAnswer(answer string) {
	s.UserAnswer = answer
}

func (s *SimpleQuestion) Valid() bool {
	return s.RigthAnswer == s.UserAnswer
}

// вопрос, у которого несколько правильных ответов
type HardQuestion struct {
	Question
	RigthAnswers []string `json:"rigth_answers"`
	UserAnswers  []string
}

func (s *HardQuestion) SetUserAnswer(answer []string) {
	s.UserAnswers = answer
}

func (s *HardQuestion) Valid() bool {
	return reflect.DeepEqual(s.UserAnswers, s.RigthAnswers)
}
