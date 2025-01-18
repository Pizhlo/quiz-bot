package controller

import (
	"fmt"
	"quiz-mod/internal/view"

	"gopkg.in/telebot.v3"
)

func (c *Controller) Answer(telectx telebot.Context) error {
	rigthAnswer, err := c.questionSrv.RigthAnswer(telectx.Chat().ID)
	if err != nil {
		return err
	}

	text := telectx.Text()

	msg := fmt.Sprintf("%s\n\nПравильный ответ: %+v", text, rigthAnswer)

	return telectx.EditOrSend(msg, view.Next())
}

func (c *Controller) Next(telectx telebot.Context) error {
	c.questionSrv.SetNext(telectx.Chat().ID)

	return c.question(telectx)
}
