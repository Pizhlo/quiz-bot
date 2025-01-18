package server

import (
	"context"
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

func (s *Server) buttonHandlers(ctx context.Context) {
	s.bot.Handle(&view.BtnStartQuiz, func(ctx telebot.Context) error {
		return ctx.EditOrSend(message.FirstMessage, view.StartFirstLevel())
	})

	s.bot.Handle(&view.BtnStartFirstLevel, func(telectx telebot.Context) error {
		return s.controller.StartFirstLevel(ctx, telectx)
	})

	s.bot.Handle(&view.BtnBackToMenu, func(ctx telebot.Context) error {
		return ctx.EditOrSend(message.StartMessage, view.MainMenu())
	})

	s.bot.Handle(&view.BtnAnswer, func(ctx telebot.Context) error {
		return s.controller.Answer(ctx)
	})

	s.bot.Handle(&view.BtnNext, func(ctx telebot.Context) error {
		return s.controller.Next(ctx)
	})
}
