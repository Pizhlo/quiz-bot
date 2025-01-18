package server

import (
	"quiz-mod/internal/controller"

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

func (s *Server) Start() {
	s.commands()

	logrus.Info("server started")
}
