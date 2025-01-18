package model

import "reflect"

type question struct {
	Text    string   `json:"question"`
	Answers []string `json:"answers"`
}

// вопрос, у которого только один правильный ответ
type SimpleQuestion struct {
	question
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
	question
	RigthAnswers []string `json:"rigth_answers"`
	UserAnswers  []string
}

func (s *HardQuestion) SetUserAnswer(answer []string) {
	s.UserAnswers = answer
}

func (s *HardQuestion) Valid() bool {
	return reflect.DeepEqual(s.UserAnswers, s.RigthAnswers)
}
