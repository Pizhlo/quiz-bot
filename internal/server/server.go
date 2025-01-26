package server

import (
	"context"
	"quiz-bot/internal/controller"

	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

type Server struct {
	bot        *tele.Bot
	controller *controller.Controller
}

func New(bot *tele.Bot, controller *controller.Controller) *Server {
	return &Server{
		bot:        bot,
		controller: controller,
	}
}

func (s *Server) Start(ctx context.Context) {
	s.commandHandlers(ctx)
	s.buttonHandlers(ctx)

	logrus.Info("server started")
}
