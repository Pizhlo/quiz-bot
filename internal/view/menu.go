package view

import (
	"fmt"

	tele "gopkg.in/telebot.v3"
)

var (
	BtnStartQuiz = tele.Btn{Text: "Начать квиз", Unique: "start_quiz"}

	BtnStartFirstLevel = tele.Btn{Text: "Начать", Unique: "start_first_lvl"}

	BtnBackToMenu = tele.Btn{Text: "⬅️Меню", Unique: "menu"}

	BtnNext = tele.Btn{Text: "Дальше➡️", Unique: "next"}

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

func Next() *tele.ReplyMarkup {
	menu := &tele.ReplyMarkup{}

	menu.Inline(
		menu.Row(BtnNext),
	)

	return menu
}

func BackToMenu() *tele.ReplyMarkup {
	menu := &tele.ReplyMarkup{}

	menu.Inline(
		menu.Row(BtnBackToMenu),
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

	btns = append(btns, BtnBackToMenu)

	menu := &tele.ReplyMarkup{}

	menu.Inline(
		menu.Split(2, btns)...,
	)

	return menu
}
