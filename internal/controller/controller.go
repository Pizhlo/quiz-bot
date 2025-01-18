package controller

import (
	"fmt"
	"quiz-mod/internal/config"
	"quiz-mod/internal/message"
	"quiz-mod/internal/service/question"
	"quiz-mod/internal/view"
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

	_, channelErr := c.bot.Send(&tele.Chat{ID: int64(c.channelID)}, msg)
	if channelErr != nil {
		logrus.Errorf("Error while sending error message to channel. Error: %+v\n", channelErr)
	}
}
