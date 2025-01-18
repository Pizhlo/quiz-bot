package controller

import (
	"context"
	"quiz-mod/internal/view"

	"gopkg.in/telebot.v3"
)

func (c *Controller) StartFirstLevel(ctx context.Context, telectx telebot.Context) error {
	// начинаем первый уровень - выставляем номер уровня и вопроса
	c.questionSrv.StartFirstLvl(telectx.Chat().ID)

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
