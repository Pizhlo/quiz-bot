package model

import "time"

type Result struct {
	Duration     time.Duration
	RigthAnswers map[int]int
}

func (s *Result) SaveAnswers(question int, result int) {
	if s.RigthAnswers == nil {
		s.RigthAnswers = make(map[int]int)
	}

	s.RigthAnswers[question] = result
}
