package server

import (
	"quiz-mod/internal/command"
	"quiz-mod/internal/message"
	"quiz-mod/internal/view"

	"gopkg.in/telebot.v3"
)

func (s *Server) commandHandlers() {
	s.bot.Handle(command.Start, func(ctx telebot.Context) error {
		return ctx.EditOrSend(message.StartMessage, view.MainMenu())
	})

	s.bot.Handle(command.Help, func(ctx telebot.Context) error {
		return ctx.EditOrSend(message.HelpMessage)
	})
}

func (s *Server) buttonHandlers() {
	s.bot.Handle(&view.BtnStartQuiz, func(ctx telebot.Context) error {
		return ctx.EditOrSend(message.FirstMessage, view.StartFirstLevel())
	})
}
