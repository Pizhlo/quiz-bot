package controller

import (
	"context"
	"quiz-bot/internal/config"
	"quiz-bot/internal/message"
	"quiz-bot/internal/service/question"
	"quiz-bot/internal/storage/postgres/quiz"
	"quiz-bot/internal/view"
	"quiz-bot/mocks"
	"quiz-bot/pkg/random"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/telebot.v3"
)

func TestReset(t *testing.T) {
	cfg := config.Config{
		FirstLevel:  random.SimpleQuestions(1, false),
		SecondLevel: random.HardQuestions(5, false),
		ThirdLevel:  random.SimpleQuestions(5, false),
	}

	qSrv := question.New(&cfg, nil, nil, "")

	// сохраняем пользователя
	qSrv.StartFirstLvl(1)

	controller := New(nil, 0, &cfg, qSrv)

	ctrl := gomock.NewController(t)
	telectx := mocks.NewMockteleCtx(ctrl)

	chat := telebot.Chat{
		ID: 1,
	}

	telectx.EXPECT().Chat().Return(&chat)
	telectx.EXPECT().Message().Return(&telebot.Message{})

	telectx.EXPECT().EditOrSend(gomock.Any(), gomock.Any()).Return(nil).Do(func(msg string, kb *telebot.ReplyMarkup) {
		expectedMsg := message.StartMessage

		expectedKb := view.MainMenu()

		assert.Equal(t, expectedMsg, msg)
		assert.Equal(t, expectedKb, kb)
	})

	err := controller.Reset(telectx)
	require.NoError(t, err)
}

func TestReset_WithPicture(t *testing.T) {
	cfg := config.Config{
		FirstLevel:  random.SimpleQuestions(1, false),
		SecondLevel: random.HardQuestions(5, false),
		ThirdLevel:  random.SimpleQuestions(5, false),
	}

	qSrv := question.New(&cfg, nil, nil, "")

	// сохраняем пользователя
	qSrv.StartFirstLvl(1)

	controller := New(nil, 0, &cfg, qSrv)

	ctrl := gomock.NewController(t)
	telectx := mocks.NewMockteleCtx(ctrl)

	chat := telebot.Chat{
		ID: 1,
	}

	telectx.EXPECT().Chat().Return(&chat)
	telectx.EXPECT().Message().Return(&telebot.Message{Caption: "test"})
	telectx.EXPECT().Delete().Return(nil)

	telectx.EXPECT().Send(gomock.Any(), gomock.Any()).Return(nil).Do(func(msg string, kb *telebot.ReplyMarkup) {
		expectedMsg := message.StartMessage

		expectedKb := view.MainMenu()

		assert.Equal(t, expectedMsg, msg)
		assert.Equal(t, expectedKb, kb)
	})

	err := controller.Reset(telectx)
	require.NoError(t, err)
}

func TestResultsByUserID_NoResults(t *testing.T) {
	cfg := config.Config{
		FirstLevel:  random.SimpleQuestions(1, false),
		SecondLevel: random.HardQuestions(5, false),
		ThirdLevel:  random.SimpleQuestions(5, false),
	}

	ctrl := gomock.NewController(t)

	telectx := mocks.NewMockteleCtx(ctrl)
	db := mocks.NewMockstorage(ctrl)

	qSrv := question.New(&cfg, db, nil, "")

	// сохраняем пользователя
	qSrv.StartFirstLvl(1)

	controller := New(nil, 0, &cfg, qSrv)

	chat := telebot.Chat{
		ID: 1,
	}

	telectx.EXPECT().Chat().Return(&chat)

	db.EXPECT().AllResults(gomock.Any(), gomock.Any()).Do(func(ctx any, userID int64) {
		assert.Equal(t, chat.ID, userID)
	}).Return(nil, quiz.ErrNoResults)

	telectx.EXPECT().EditOrSend(gomock.Any(), gomock.Any()).Return(nil).Do(func(msg string, kb *telebot.ReplyMarkup) {
		expectedMsg := message.NoResultsMessage

		expectedKb := view.BackToMenu()

		assert.Equal(t, expectedMsg, msg)
		assert.Equal(t, expectedKb, kb)
	})

	err := controller.ResultsByUserID(context.TODO(), telectx)
	require.NoError(t, err)
}

func TestResultsByUserID(t *testing.T) {
	cfg := config.Config{
		FirstLevel:  random.SimpleQuestions(1, false),
		SecondLevel: random.HardQuestions(5, false),
		ThirdLevel:  random.SimpleQuestions(5, false),
	}

	ctrl := gomock.NewController(t)

	telectx := mocks.NewMockteleCtx(ctrl)
	db := mocks.NewMockstorage(ctrl)

	qSrv := question.New(&cfg, db, nil, "")

	// сохраняем пользователя
	qSrv.StartFirstLvl(1)

	controller := New(nil, 0, &cfg, qSrv)

	chat := telebot.Chat{
		ID: 1,
	}

	dbResults := random.Results(10)

	telectx.EXPECT().Chat().Return(&chat)

	db.EXPECT().AllResults(gomock.Any(), gomock.Any()).Do(func(ctx any, userID int64) {
		assert.Equal(t, chat.ID, userID)
	}).Return(dbResults, nil)

	telectx.EXPECT().EditOrSend(gomock.Any(), gomock.Any()).Return(nil).Do(func(msg string, kb *telebot.ReplyMarkup) {
		view := view.New()

		expectedMsg := view.Message(dbResults)

		expectedKb := view.Keyboard()

		assert.Equal(t, expectedMsg, msg)
		assert.Equal(t, expectedKb, kb)
	})

	err := controller.ResultsByUserID(context.TODO(), telectx)
	require.NoError(t, err)
}
