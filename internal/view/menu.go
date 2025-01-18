package view

import (
	"fmt"

	tele "gopkg.in/telebot.v3"
)

var (
	BtnStartQuiz = tele.Btn{Text: "Начать квиз", Unique: "start_quiz"}

	BtnStartFirstLevel = tele.Btn{Text: "Начать", Unique: "start_first_lvl"}

	// кнопка для ответов
	BtnAnswer = tele.Btn{Unique: "answer"}
)

func MainMenu() *tele.ReplyMarkup {
	menu := &tele.ReplyMarkup{}

	menu.Inline(
		menu.Row(BtnStartQuiz),
	)

	return menu
}

func StartFirstLevel() *tele.ReplyMarkup {
	menu := &tele.ReplyMarkup{}

	menu.Inline(
		menu.Row(BtnStartFirstLevel),
	)

	return menu
}

func Answers(answers []string) *tele.ReplyMarkup {
	btns := []tele.Btn{}

	for i, answer := range answers {
		BtnAnswer.Text = answer
		BtnAnswer.Data = fmt.Sprintf("%d", i)

		btns = append(btns, BtnAnswer)
	}

	menu := &tele.ReplyMarkup{}

	menu.Inline(
		menu.Row(btns...),
	)

	return menu
}
