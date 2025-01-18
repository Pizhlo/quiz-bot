package controller

import (
	"context"

	"gopkg.in/telebot.v3"
)

func (c *Controller) StartFirstLevel(ctx context.Context, telectx telebot.Context) error {
	// начинаем первый уровень - выставляем номер уровня и вопроса
	c.questionSrv.StartFirstLvl(telectx.Chat().ID)

	return c.sendCurrentQuestion(telectx)
}

func (c *Controller) StartSecondLevel(ctx context.Context, telectx telebot.Context) error {
	// начинаем первый уровень - выставляем номер уровня и вопроса
	err := c.questionSrv.StartSecondLvl(telectx.Chat().ID)
	if err != nil {
		return err
	}

	return c.sendCurrentQuestion(telectx)
}

func (c *Controller) StartThirdLevel(ctx context.Context, telectx telebot.Context) error {
	// начинаем первый уровень - выставляем номер уровня и вопроса
	err := c.questionSrv.StartThirdLvl(telectx.Chat().ID)
	if err != nil {
		return err
	}

	return c.sendCurrentQuestion(telectx)
}
