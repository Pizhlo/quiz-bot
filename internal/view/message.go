package view

import (
	"fmt"
	"quiz-mod/internal/message"
	"quiz-mod/internal/model"
)

const (
	dateFieldFormat = "02.01.2006 15:04:05"
)

type ResultView struct {
	pages       []string
	currentPage int
}

func NewResult() *ResultView {
	return &ResultView{pages: make([]string, 0), currentPage: 0}
}

func (s *ResultView) Message(results []model.Result) string {
	page := "<b>Твои результаты:\n\n</b>"

	for i, res := range results {
		msg := ""

		// seconds := time.Duration(res.Seconds) * time.Second
		msg += fmt.Sprintf(message.Result,
			res.RigthAnswers[model.FirstLevel],
			res.TotalAnswers[model.FirstLevel],
			res.RigthAnswers[model.SecondLevel],
			res.TotalAnswers[model.SecondLevel],
			res.RigthAnswers[model.ThirdLevel],
			res.TotalAnswers[model.ThirdLevel],
			// res.Duration.String(),
			fmt.Sprintf("%.2fs", res.Seconds),
			// seconds.String(),
		)

		page += fmt.Sprintf("%s\n\nДата: %s", msg, res.Date.Format(dateFieldFormat))

		if i%5 == 0 {
			s.pages = append(s.pages, page)
		}

		page = ""
	}

	return s.pages[0]
}
