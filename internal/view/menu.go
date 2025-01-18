package view

import (
	tele "gopkg.in/telebot.v3"
)

var (
	BtnStartQuiz = tele.Btn{Text: "Начать квиз", Unique: "start_quiz"}
	BtnResults   = tele.Btn{Text: "Мои результаты", Unique: "results"}

	BtnNewLvl = tele.Btn{Text: "Дальше➡️", Unique: "new_lvl"}

	BtnStartFirstLevel  = tele.Btn{Text: "Начать", Unique: "start_first_lvl"}
	BtnStartSecondLevel = tele.Btn{Text: "Начать", Unique: "start_second_lvl"}
	BtnStartThirdLevel  = tele.Btn{Text: "Начать", Unique: "start_third_lvl"}

	BtnBackToMenu = tele.Btn{Text: "⬅️Меню", Unique: "menu"}

	BtnNext = tele.Btn{Text: "Дальше➡️", Unique: "next"}

	// кнопка для простых ответов ответов
	BtnSimpleAnswer = tele.Btn{Unique: "simple_answer"}

	// кнопка для 2 раунда, где несколько вариантов ответа
	BtnAnswer = tele.Btn{Unique: "answer"}
	// кнопка для отправки ответа (для множественного выбора)
	BtnSendAnswer = tele.Btn{Text: "📝Ответить", Unique: "send_answer"}
)

func MainMenu() *tele.ReplyMarkup {
	menu := &tele.ReplyMarkup{}

	menu.Inline(
		menu.Row(BtnStartQuiz),
		menu.Row(BtnResults),
	)

	return menu
}

func Next() *tele.ReplyMarkup {
	menu := &tele.ReplyMarkup{}

	menu.Inline(
		menu.Row(BtnNext),
		menu.Row(BtnBackToMenu),
	)

	return menu
}

func NewLvl() *tele.ReplyMarkup {
	menu := &tele.ReplyMarkup{}

	menu.Inline(
		menu.Row(BtnNewLvl),
		menu.Row(BtnBackToMenu),
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
		menu.Row(BtnBackToMenu),
	)

	return menu
}

func StartSecondLevel() *tele.ReplyMarkup {
	menu := &tele.ReplyMarkup{}

	menu.Inline(
		menu.Row(BtnStartSecondLevel),
		menu.Row(BtnBackToMenu),
	)

	return menu
}

func StartThirdLevel() *tele.ReplyMarkup {
	menu := &tele.ReplyMarkup{}

	menu.Inline(
		menu.Row(BtnStartThirdLevel),
		menu.Row(BtnBackToMenu),
	)

	return menu
}

func SimpleAnswers(answers []string) *tele.ReplyMarkup {
	btns := []tele.Btn{}

	for _, answer := range answers {
		BtnSimpleAnswer.Text = answer
		BtnSimpleAnswer.Data = answer

		btns = append(btns, BtnSimpleAnswer)
	}

	btns = append(btns, BtnBackToMenu)

	menu := &tele.ReplyMarkup{}

	menu.Inline(
		menu.Split(2, btns)...,
	)

	return menu
}

func Answers(answers []string) *tele.ReplyMarkup {
	btns := []tele.Btn{}

	for _, answer := range answers {
		BtnAnswer.Text = answer
		BtnAnswer.Data = answer

		btns = append(btns, BtnAnswer)
	}

	btns = append(btns, BtnBackToMenu, BtnSendAnswer)

	menu := &tele.ReplyMarkup{}

	menu.Inline(
		menu.Split(2, btns)...,
	)

	return menu
}
