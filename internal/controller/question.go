package controller

import (
	"quiz-mod/internal/view"

	"gopkg.in/telebot.v3"
)

func (c *Controller) sendCurrentQuestion(telectx telebot.Context) error {
	question, err := c.questionSrv.CurrentQuestion(telectx.Chat().ID)
	if err != nil {
		return err
	}

	msg, err := c.questionSrv.Message(telectx.Chat().ID)
	if err != nil {
		return err
	}

	btns := view.Answers(question.Answers)

	return telectx.EditOrSend(msg, btns)
}
