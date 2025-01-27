package controller

import (
	"context"
	"quiz-bot/internal/config"
	"quiz-bot/internal/service/question"
	"quiz-bot/internal/view"
	"quiz-bot/mocks"
	"quiz-bot/pkg/random"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/telebot.v3"
)

func TestStartFirstLevel(t *testing.T) {
	cfg := config.Config{
		FirstLevel:  random.SimpleQuestions(5, false),
		SecondLevel: random.HardQuestions(5, false),
		ThirdLevel:  random.SimpleQuestions(5, false),
	}
	qSrv := question.New(&cfg, nil, nil)

	controller := New(nil, 0, &cfg, qSrv)

	ctrl := gomock.NewController(t)
	telectx := mocks.NewMockteleCtx(ctrl)

	telectx.EXPECT().EditOrSend(gomock.Any(), gomock.Any()).Return(nil).Do(func(msg string, kb *telebot.ReplyMarkup) {
		expectedMsg, err := qSrv.Message(1)
		require.NoError(t, err)

		expectedKb := view.SimpleAnswers(cfg.FirstLevel[0].Answers)

		assert.Equal(t, expectedMsg, msg)
		assert.Equal(t, expectedKb, kb)
	})

	chat := telebot.Chat{
		ID: 1,
	}

	telectx.EXPECT().Chat().Return(&chat).Times(4)

	err := controller.StartFirstLevel(context.TODO(), telectx)
	require.NoError(t, err)
}
