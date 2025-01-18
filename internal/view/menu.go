package view

import tele "gopkg.in/telebot.v3"

var (
	BtnStartQuiz = tele.Btn{Text: "Начать квиз", Unique: "start_quiz"}
)

func MainMenu() *tele.ReplyMarkup {
	menu := &tele.ReplyMarkup{}

	menu.Inline(
		menu.Row(BtnStartQuiz),
	)

	return menu
}
