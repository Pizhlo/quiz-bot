package controller

import (
	"fmt"
	"quiz-bot/internal/config"
	"quiz-bot/internal/message"
	"quiz-bot/internal/service/question"
	"quiz-bot/internal/view"
	"sync"

	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

type Controller struct {
	mu          sync.Mutex
	bot         *tele.Bot
	channelID   int
	cfg         *config.Config
	questionSrv *question.Question
}

func New(bot *tele.Bot,
	channelID int,
	cfg *config.Config,
	questionSrv *question.Question) *Controller {

	return &Controller{
		bot:         bot,
		mu:          sync.Mutex{},
		cfg:         cfg,
		channelID:   channelID,
		questionSrv: questionSrv,
	}
}

//lint:ignore U1000 Ignore unused function temporarily for debugging
//go:generate mockgen -source ./controller.go -destination=../../mocks/controller.go -package=mocks
type teleCtx interface {
	tele.Context
}

// HandleError сообщает об ошибке в канал.
// Также сообщает пользователю об ошибке - редактирует сообщение
func (c *Controller) HandleError(ctx tele.Context, err error) {
	var msg string
	if ctx.Message().Sender != c.bot.Me {
		msg = fmt.Sprintf(message.ErrorMessageChannelMessage, ctx.Message().Text, err)
	} else {
		msg = fmt.Sprintf(message.ErrorMessageChannelMessage, ctx.Callback().Unique, err)
	}

	editErr := ctx.EditOrSend(message.ErrorMessageUser, view.BackToMenu())
	if editErr != nil {
		logrus.Errorf("Error while sending error message to user. Error: %+v\n", editErr)
	}

	logrus.Debug(msg)

	_, channelErr := c.bot.Send(&tele.Chat{ID: int64(c.channelID)}, msg)
	if channelErr != nil {
		logrus.Errorf("Error while sending error message to channel. Error: %+v\n", channelErr)
	}
}
