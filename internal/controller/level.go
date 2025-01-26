package controller

import (
	"context"
	"fmt"
	"quiz-bot/internal/message"
	"quiz-bot/internal/model"
	"quiz-bot/internal/view"

	"github.com/sirupsen/logrus"
	"gopkg.in/telebot.v3"
)

func (c *Controller) StartFirstLevel(ctx context.Context, telectx telebot.Context) error {
	// начинаем первый уровень - выставляем номер уровня и вопроса
	c.questionSrv.StartFirstLvl(telectx.Chat().ID)

	return c.sendCurrentQuestion(ctx, telectx)
}

func (c *Controller) StartSecondLevel(ctx context.Context, telectx telebot.Context) error {
	// начинаем первый уровень - выставляем номер уровня и вопроса
	err := c.questionSrv.StartSecondLvl(telectx.Chat().ID)
	if err != nil {
		return err
	}

	return c.sendCurrentQuestion(ctx, telectx)
}

func (c *Controller) StartThirdLevel(ctx context.Context, telectx telebot.Context) error {
	// начинаем первый уровень - выставляем номер уровня и вопроса
	err := c.questionSrv.StartThirdLvl(telectx.Chat().ID)
	if err != nil {
		return err
	}

	return c.sendCurrentQuestion(ctx, telectx)
}

func (c *Controller) Next(ctx context.Context, telectx telebot.Context) error {
	last, err := c.questionSrv.IsQuestionLast(telectx.Chat().ID)
	if err != nil {
		return err
	}

	// если вопрос не последний - отправляем следующий вопрос
	if !last {
		c.questionSrv.SetNext(telectx.Chat().ID)

		return c.sendCurrentQuestion(ctx, telectx)
	}

	// отправляем сообщение с результатами раунда
	return c.levelResuls(telectx)
}

// sendLevelMessage отправляет сообщение с описанием уровня
func (c *Controller) SendLevelMessage(ctx context.Context, telectx telebot.Context) error {
	lvl, err := c.questionSrv.CurrentLevel(telectx.Chat().ID)
	if err != nil {
		return err
	}

	switch lvl {
	case model.FirstLevel:
		c.questionSrv.SetNext(telectx.Chat().ID)

		if telectx.Message().Caption != "" {
			err := telectx.Delete()
			if err != nil {
				logrus.Errorf("error deleting message: %+v", err)
			}

			return telectx.Send(message.SecondLvlMessage, view.StartSecondLevel())
		}

		return telectx.EditOrSend(message.SecondLvlMessage, view.StartSecondLevel())
	case model.SecondLevel:
		c.questionSrv.SetNext(telectx.Chat().ID)
		if telectx.Message().Caption != "" {
			err := telectx.Delete()
			if err != nil {
				logrus.Errorf("error deleting message: %+v", err)
			}

			return telectx.Send(message.ThirdLvlMessage, view.StartThirdLevel())
		}

		return telectx.EditOrSend(message.ThirdLvlMessage, view.StartThirdLevel())
	case model.ThirdLevel:
		if telectx.Message().Caption != "" {
			err := telectx.Delete()
			if err != nil {
				logrus.Errorf("error deleting message: %+v", err)
			}
		}

		return c.results(ctx, telectx)
	default:
		return fmt.Errorf("unknown lvl: %+v", lvl)
	}
}
