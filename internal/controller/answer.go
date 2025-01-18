package controller

import (
	"fmt"
	"quiz-mod/internal/message"
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
	last, err := c.questionSrv.IsQuestionLast(telectx.Chat().ID)
	if err != nil {
		return err
	}

	// если вопрос не последний - отправляем следующий вопрос
	if !last {
		c.questionSrv.SetNext(telectx.Chat().ID)

		return c.sendCurrentQuestion(telectx)
	}

	// отправляем сообщение с описанием следующего раунда
	return c.sendLevelMessage(telectx)
}

func (c *Controller) sendLevelMessage(telectx telebot.Context) error {
	lvl, err := c.questionSrv.CurrentLevel(telectx.Chat().ID)
	if err != nil {
		return err
	}

	switch lvl {
	case 0:
		c.questionSrv.SetNext(telectx.Chat().ID)
		return telectx.EditOrSend(message.SecondLvlMessage, view.StartSecondLevel())
	case 1:
		c.questionSrv.SetNext(telectx.Chat().ID)
		return telectx.EditOrSend(message.SecondLvlMessage, view.StartSecondLevel())
	case 2:
		return c.results(telectx)
	default:
		return fmt.Errorf("unknown lvl: %+v", lvl)
	}
}

func (c *Controller) results(telectx telebot.Context) error {
	return telectx.EditOrSend("results", view.BackToMenu())
}
