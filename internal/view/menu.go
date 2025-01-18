package view

import tele "gopkg.in/telebot.v3"

var (
	BtnStartQuiz = tele.Btn{Text: "Начать квиз", Unique: "start_quiz"}

	BtnStartFirstLevel = tele.Btn{Text: "Начать", Unique: "start_first_lvl"}
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
