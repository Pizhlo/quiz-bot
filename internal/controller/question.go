package controller

import (
	"fmt"
	"quiz-mod/internal/view"

	"gopkg.in/telebot.v3"
)

func (c *Controller) sendCurrentQuestion(telectx telebot.Context) error {
	lvl, err := c.questionSrv.CurrentLevel(telectx.Chat().ID)
	if err != nil {
		return err
	}

	question, err := c.questionSrv.CurrentQuestion(telectx.Chat().ID)
	if err != nil {
		return err
	}

	msg, err := c.questionSrv.Message(telectx.Chat().ID)
	if err != nil {
		return err
	}

	var btns *telebot.ReplyMarkup

	switch lvl {
	case 0, 2:
		btns = view.SimpleAnswers(question.Answers)
	case 1:
		btns = view.Answers(question.Answers)
	default:
		return fmt.Errorf("invalid level: %+v", lvl)
	}

	return telectx.EditOrSend(msg, btns)
}
