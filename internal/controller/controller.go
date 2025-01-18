package controller

import (
	"quiz-mod/internal/service/question"
	"sync"

	tele "gopkg.in/telebot.v3"
)

type Controller struct {
	mu          sync.Mutex
	bot         *tele.Bot
	channelID   int
	questionSrv *question.Question
}

func New(bot *tele.Bot,
	channelID int,
	questionSrv *question.Question) *Controller {

	return &Controller{
		bot:         bot,
		mu:          sync.Mutex{},
		channelID:   channelID,
		questionSrv: questionSrv,
	}
}
