package controller

import (
	"sync"

	tele "gopkg.in/telebot.v3"
)

type Controller struct {
	mu        sync.Mutex
	bot       *tele.Bot
	channelID int
}

func New(bot *tele.Bot,
	channelID int) *Controller {

	return &Controller{
		bot:       bot,
		mu:        sync.Mutex{},
		channelID: channelID,
	}
}
