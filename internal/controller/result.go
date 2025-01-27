package controller

import (
	"context"
	"errors"
	"fmt"
	"quiz-bot/internal/message"
	"quiz-bot/internal/model"
	"quiz-bot/internal/storage/postgres/quiz"
	"quiz-bot/internal/view"

	"github.com/sirupsen/logrus"
	"gopkg.in/telebot.v3"
)

func (c *Controller) results(ctx context.Context, telectx telebot.Context) error {
	// засекаем, сколько времени прошло с начала викторины
	err := c.questionSrv.StopTimer(telectx.Chat().ID)
	if err != nil {
		return err
	}

	// сохраняем результаты в БД
	err = c.questionSrv.SaveResults(ctx, telectx.Chat().ID)
	if err != nil {
		return err
	}

	// получаем результаты, чтобы показать пользователю
	res, err := c.questionSrv.Results(telectx.Chat().ID)
	if err != nil {
		return err
	}

	// отправляем результаты пользователю
	err = c.sendResultsToUser(telectx, res)
	if err != nil {
		return err
	}

	username := ""

	if telectx.Chat().Username != "" {
		username = fmt.Sprintf("@%s", telectx.Chat().Username)
	} else {
		username = fmt.Sprintf("%s %s", telectx.Chat().FirstName, telectx.Chat().LastName)
	}

	// отправляем результаты в канал
	return c.sendResultsToChan(username, res)
}

func (c *Controller) sendResultsToUser(telectx telebot.Context, res model.Result) error {
	result := c.resultMsg(res)

	msg := fmt.Sprintf(message.ResultMessage, result)

	return telectx.EditOrSend(msg, view.ResultMenu())
}

func (c *Controller) resultMsg(res model.Result) string {
	return fmt.Sprintf(message.Result, res.RigthAnswers[0], len(c.cfg.FirstLevel),
		res.RigthAnswers[1], len(c.cfg.SecondLevel),
		res.RigthAnswers[2], len(c.cfg.ThirdLevel),
		fmt.Sprintf("%.2fs", res.Seconds))
}

func (c *Controller) sendResultsToChan(username string, res model.Result) error {
	result := c.resultMsg(res)

	msg := fmt.Sprintf(message.ChannelResultMessage, username, result)

	_, err := c.bot.Send(&telebot.Chat{ID: int64(c.channelID)}, msg)
	return err
}

func (c *Controller) levelResuls(telectx telebot.Context) error {
	rigthAns, err := c.questionSrv.LevelResults(telectx.Chat().ID)
	if err != nil {
		return err
	}

	questionsNum, err := c.questionSrv.QuestionNum(telectx.Chat().ID)
	if err != nil {
		return err
	}

	msg := fmt.Sprintf(message.LevelEnd, rigthAns, questionsNum)

	if telectx.Message().Caption != "" {
		err := telectx.Delete()
		if err != nil {
			logrus.Errorf("error deleting message: %+v", err)
		}

		return telectx.Send(msg, view.NewLvl())
	}

	return telectx.EditOrSend(msg, view.NewLvl())
}

// Reset сбрасывает все сохраненные данные. Используется, если пользователь ушел в главное меню
func (c *Controller) Reset(telectx telebot.Context) error {
	c.questionSrv.Reset(telectx.Chat().ID)

	if telectx.Message().Caption != "" {
		err := telectx.Delete()
		if err != nil {
			logrus.Errorf("error deleting message: %+v", err)
		}

		return telectx.Send(message.StartMessage, view.MainMenu())
	}

	return telectx.EditOrSend(message.StartMessage, view.MainMenu())
}

// ResultsByUserID обрабатывает нажатие на кнопку "Мои результаты".
// Достает из БД все результаты викторин пользователя и отправляет сообщение
func (c *Controller) ResultsByUserID(ctx context.Context, telectx telebot.Context) error {
	msg, kb, err := c.questionSrv.AllResults(ctx, telectx.Chat().ID)
	if err != nil {
		if errors.Is(err, quiz.ErrNoResults) {
			return telectx.EditOrSend(message.NoResultsMessage, view.BackToMenu())
		}

		return err
	}

	return telectx.EditOrSend(msg, kb)
}
