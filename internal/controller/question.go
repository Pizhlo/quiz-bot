package controller

import (
	"context"
	"fmt"
	"os"
	"quiz-bot/internal/model"
	"quiz-bot/internal/view"

	"github.com/sirupsen/logrus"
	"gopkg.in/telebot.v3"
)

func (c *Controller) sendCurrentQuestion(ctx context.Context, telectx telebot.Context) error {
	lvl, err := c.questionSrv.CurrentLevel(telectx.Chat().ID)
	if err != nil {
		return err
	}

	question, err := c.questionSrv.CurrentQuestion(telectx.Chat().ID)
	if err != nil {
		return err
	}

	if question.Picture != "" {
		kb, err := keyboardFromQuestion(question, lvl)
		if err != nil {
			return err
		}

		msg, err := c.questionSrv.Message(telectx.Chat().ID)
		if err != nil {
			return err
		}

		return c.sendPicture(ctx, telectx, msg, question.Picture, kb)
	}

	return c.sendTextWithBtns(telectx, question, lvl)
}

func (c *Controller) sendPicture(ctx context.Context, telectx telebot.Context, msg string, filename string, kb *telebot.ReplyMarkup) error {
	file, err := c.getPicture(ctx, filename)
	if err != nil {
		return err
	}

	return telectx.EditOrSend(&telebot.Photo{File: telebot.FromReader(file), Caption: msg}, &telebot.SendOptions{
		ReplyMarkup: kb,
		ParseMode:   htmlParseMode,
	})
}

func (c *Controller) getPicture(ctx context.Context, filename string) (*os.File, error) {
	err := c.questionSrv.GetFile(ctx, filename)
	if err != nil {
		return nil, err
	}

	return os.Open(filename)
}

func keyboardFromQuestion(question *model.Question, lvl int) (*telebot.ReplyMarkup, error) {
	switch lvl {
	case model.FirstLevel, model.ThirdLevel:
		return view.SimpleAnswers(question.Answers), nil
	case model.SecondLevel:
		return view.Answers(question.Answers), nil
	default:
		return nil, fmt.Errorf("invalid level: %+v", lvl)
	}
}

func (c *Controller) sendTextWithBtns(telectx telebot.Context, question *model.Question, lvl int) error {
	msg, err := c.questionSrv.Message(telectx.Chat().ID)
	if err != nil {
		return err
	}

	kb, err := keyboardFromQuestion(question, lvl)
	if err != nil {
		return err
	}

	if telectx.Message().Caption != "" {
		err := telectx.Delete()
		if err != nil {
			logrus.Errorf("error deleting message: %+v", err)
		}

		return telectx.Send(msg, kb)
	}

	return telectx.EditOrSend(msg, kb)
}
