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
		err := ctx.EditOrSend(message.StartMessage, view.MainMenu())
		if err != nil {
			s.controller.HandleError(ctx, err)
		}

		return nil
	})

	s.bot.Handle(command.Help, func(ctx telebot.Context) error {
		err := ctx.EditOrSend(message.HelpMessage)
		if err != nil {
			s.controller.HandleError(ctx, err)
		}

		return nil
	})

	s.bot.Handle(telebot.OnText, func(ctx telebot.Context) error {
		err := s.controller.OnText(ctx)
		if err != nil {
			s.controller.HandleError(ctx, err)
		}

		return nil
	})
}

func (s *Server) buttonHandlers(ctx context.Context) {
	// main menu
	s.bot.Handle(&view.BtnStartQuiz, func(ctx telebot.Context) error {
		err := ctx.EditOrSend(message.FirstLvlMessage, view.StartFirstLevel())
		if err != nil {
			s.controller.HandleError(ctx, err)
		}

		return nil
	})

	s.bot.Handle(&view.BtnResults, func(telectx telebot.Context) error {
		err := s.controller.ResultsByUserID(ctx, telectx)
		if err != nil {
			s.controller.HandleError(telectx, err)
		}

		return nil
	})

	s.bot.Handle(&view.BtnBackToMenu, func(ctx telebot.Context) error {
		err := s.controller.Reset(ctx)
		if err != nil {
			s.controller.HandleError(ctx, err)
		}

		return nil
	})

	// lvls
	s.bot.Handle(&view.BtnStartFirstLevel, func(telectx telebot.Context) error {
		err := s.controller.StartFirstLevel(ctx, telectx)
		if err != nil {
			s.controller.HandleError(telectx, err)
		}

		return nil
	})

	s.bot.Handle(&view.BtnStartSecondLevel, func(telectx telebot.Context) error {
		err := s.controller.StartSecondLevel(ctx, telectx)
		if err != nil {
			s.controller.HandleError(telectx, err)
		}

		return nil
	})

	s.bot.Handle(&view.BtnStartThirdLevel, func(telectx telebot.Context) error {
		err := s.controller.StartThirdLevel(ctx, telectx)
		if err != nil {
			s.controller.HandleError(telectx, err)
		}

		return nil
	})

	// answers
	s.bot.Handle(&view.BtnSimpleAnswer, func(ctx telebot.Context) error {
		err := s.controller.SimpleAnswer(ctx)
		if err != nil {
			s.controller.HandleError(ctx, err)
		}

		return nil
	})

	s.bot.Handle(&view.BtnAnswer, func(ctx telebot.Context) error {
		err := s.controller.MultipleAnswer(ctx)
		if err != nil {
			s.controller.HandleError(ctx, err)
		}

		return nil
	})

	s.bot.Handle(&view.BtnSendAnswer, func(ctx telebot.Context) error {
		err := s.controller.SendAnswer(ctx)
		if err != nil {
			s.controller.HandleError(ctx, err)
		}

		return nil
	})

	s.bot.Handle(&view.BtnNext, func(ctx telebot.Context) error {
		err := s.controller.Next(ctx)
		if err != nil {
			s.controller.HandleError(ctx, err)
		}

		return nil
	})

	s.bot.Handle(&view.BtnNewLvl, func(telectx telebot.Context) error {
		err := s.controller.SendLevelMessage(ctx, telectx)
		if err != nil {
			s.controller.HandleError(telectx, err)
		}

		return nil
	})

	// pages

	s.bot.Handle(&view.BtnPrevPgResults, func(telectx telebot.Context) error {
		err := s.controller.PrevPage(telectx)
		if err != nil {
			s.controller.HandleError(telectx, err)
		}

		return nil
	})

	s.bot.Handle(&view.BtnNextPgResults, func(telectx telebot.Context) error {
		err := s.controller.NextPage(telectx)
		if err != nil {
			s.controller.HandleError(telectx, err)
		}

		return nil
	})

	s.bot.Handle(&view.BtnFirstPgResults, func(telectx telebot.Context) error {
		err := s.controller.FirstPage(telectx)
		if err != nil {
			s.controller.HandleError(telectx, err)
		}

		return nil
	})

	s.bot.Handle(&view.BtnLastPgResults, func(telectx telebot.Context) error {
		err := s.controller.LastPage(telectx)
		if err != nil {
			s.controller.HandleError(telectx, err)
		}

		return nil
	})
}
