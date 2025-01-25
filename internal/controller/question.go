package controller

import (
	"context"
	"fmt"
	"os"
	"quiz-mod/internal/model"
	"quiz-mod/internal/view"

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
		return c.sendPicture(ctx, telectx, question, lvl)
	}

	return c.sendTextWithBtns(telectx, question, lvl)
}

func (c *Controller) sendPicture(ctx context.Context, telectx telebot.Context, question *model.Question, lvl int) error {
	err := c.questionSrv.GetFile(ctx, question.Picture)
	if err != nil {
		return err
	}

	file, err := os.Open(question.Picture)
	if err != nil {
		return err
	}

	photo := &telebot.Photo{File: telebot.FromReader(file)}

	_, err = photo.Send(c.bot, telectx.Chat(), nil)
	if err != nil {
		return err
	}

	return c.sendTextWithBtns(telectx, question, lvl)
}

func (c *Controller) sendTextWithBtns(telectx telebot.Context, question *model.Question, lvl int) error {
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
