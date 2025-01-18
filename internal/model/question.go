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

func (s *SimpleQuestion) Valid(answer string) bool {
	return s.RigthAnswer == answer
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

func (s *HardQuestion) Valid(answers []string) bool {
	return reflect.DeepEqual(answers, s.RigthAnswers)
}
