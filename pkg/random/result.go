package random

import (
	"quiz-bot/internal/model"
	"time"
)

func Results(n int) []model.Result {
	res := []model.Result{}

	for i := 0; i < n; i++ {
		res = append(res, Result())
	}

	return res
}

func Result() model.Result {
	return model.Result{
		RigthAnswers: map[int]int{
			model.FirstLevel:  Int(0, 3),
			model.SecondLevel: Int(0, 3),
			model.ThirdLevel:  Int(0, 3),
		},
		TotalAnswers: map[int]int{
			model.FirstLevel:  Int(0, 3),
			model.SecondLevel: Int(0, 3),
			model.ThirdLevel:  Int(0, 3),
		},
		Date:    time.Now(),
		Seconds: float64(Int(20, 30)),
	}
}
