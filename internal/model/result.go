package model

import "time"

type Result struct {
	TgID         int64
	Duration     time.Duration
	RigthAnswers map[int]int // количество правильных ответов пользователя по раундам
	Date         time.Time
}

func (s *Result) SaveAnswers(question int, result int) {
	if s.RigthAnswers == nil {
		s.RigthAnswers = make(map[int]int)
	}

	s.RigthAnswers[question] = result
}
