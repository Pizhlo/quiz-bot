package controller

import (
	"context"
	"fmt"
	"quiz-bot/internal/config"
	"quiz-bot/internal/message"
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

func TestStartSecondLevel(t *testing.T) {
	cfg := config.Config{
		FirstLevel:  random.SimpleQuestions(5, false),
		SecondLevel: random.HardQuestions(5, false),
		ThirdLevel:  random.SimpleQuestions(5, false),
	}

	qSrv := question.New(&cfg, nil, nil)

	// сохраняем пользователя
	qSrv.StartFirstLvl(1)

	controller := New(nil, 0, &cfg, qSrv)

	ctrl := gomock.NewController(t)
	telectx := mocks.NewMockteleCtx(ctrl)

	telectx.EXPECT().EditOrSend(gomock.Any(), gomock.Any()).Return(nil).Do(func(msg string, kb *telebot.ReplyMarkup) {
		expectedMsg, err := qSrv.Message(1)
		require.NoError(t, err)

		expectedKb := view.Answers(cfg.SecondLevel[0].Answers)

		assert.Equal(t, expectedMsg, msg)
		assert.Equal(t, expectedKb, kb)
	})

	chat := telebot.Chat{
		ID: 1,
	}

	telectx.EXPECT().Chat().Return(&chat).Times(4)

	err := controller.StartSecondLevel(context.TODO(), telectx)
	require.NoError(t, err)
}

func TestStartThirdLevel(t *testing.T) {
	cfg := config.Config{
		FirstLevel:  random.SimpleQuestions(5, false),
		SecondLevel: random.HardQuestions(5, false),
		ThirdLevel:  random.SimpleQuestions(5, false),
	}

	qSrv := question.New(&cfg, nil, nil)

	// сохраняем пользователя
	qSrv.StartFirstLvl(1)

	controller := New(nil, 0, &cfg, qSrv)

	ctrl := gomock.NewController(t)
	telectx := mocks.NewMockteleCtx(ctrl)

	telectx.EXPECT().EditOrSend(gomock.Any(), gomock.Any()).Return(nil).Do(func(msg string, kb *telebot.ReplyMarkup) {
		expectedMsg, err := qSrv.Message(1)
		require.NoError(t, err)

		expectedKb := view.SimpleAnswers(cfg.ThirdLevel[0].Answers)

		assert.Equal(t, expectedMsg, msg)
		assert.Equal(t, expectedKb, kb)
	})

	chat := telebot.Chat{
		ID: 1,
	}

	telectx.EXPECT().Chat().Return(&chat).Times(4)

	err := controller.StartThirdLevel(context.TODO(), telectx)
	require.NoError(t, err)
}

func TestNext_QuestionNotLast(t *testing.T) {
	cfg := config.Config{
		FirstLevel:  random.SimpleQuestions(5, false),
		SecondLevel: random.HardQuestions(5, false),
		ThirdLevel:  random.SimpleQuestions(5, false),
	}

	qSrv := question.New(&cfg, nil, nil)

	// сохраняем пользователя
	qSrv.StartFirstLvl(1)

	controller := New(nil, 0, &cfg, qSrv)

	ctrl := gomock.NewController(t)
	telectx := mocks.NewMockteleCtx(ctrl)

	telectx.EXPECT().EditOrSend(gomock.Any(), gomock.Any()).Return(nil).Do(func(msg string, kb *telebot.ReplyMarkup) {
		expectedMsg, err := qSrv.Message(1)
		require.NoError(t, err)

		expectedKb := view.SimpleAnswers(cfg.FirstLevel[1].Answers)

		assert.Equal(t, expectedMsg, msg)
		assert.Equal(t, expectedKb, kb)
	})

	chat := telebot.Chat{
		ID: 1,
	}

	telectx.EXPECT().Chat().Return(&chat).Times(5)

	err := controller.Next(context.TODO(), telectx)
	require.NoError(t, err)
}

func TestNext_QuestionLast(t *testing.T) {
	cfg := config.Config{
		FirstLevel:  random.SimpleQuestions(1, false),
		SecondLevel: random.HardQuestions(5, false),
		ThirdLevel:  random.SimpleQuestions(5, false),
	}

	qSrv := question.New(&cfg, nil, nil)

	// сохраняем пользователя
	qSrv.StartFirstLvl(1)

	controller := New(nil, 0, &cfg, qSrv)

	ctrl := gomock.NewController(t)
	telectx := mocks.NewMockteleCtx(ctrl)

	telectx.EXPECT().EditOrSend(gomock.Any(), gomock.Any()).Return(nil).Do(func(msg string, kb *telebot.ReplyMarkup) {
		rigthAns, err := qSrv.LevelResults(telectx.Chat().ID)
		require.NoError(t, err)

		expectedMsg := fmt.Sprintf(message.LevelEnd, rigthAns, 1)

		expectedKb := view.NewLvl()

		assert.Equal(t, expectedMsg, msg)
		assert.Equal(t, expectedKb, kb)
	})

	chat := telebot.Chat{
		ID: 1,
	}

	telectx.EXPECT().Chat().Return(&chat).Times(5)
	telectx.EXPECT().Message().Return(&telebot.Message{})

	err := controller.Next(context.TODO(), telectx)
	require.NoError(t, err)
}

func TestSendLevelMessage_FirstLvl(t *testing.T) {
	cfg := config.Config{
		FirstLevel:  random.SimpleQuestions(1, false),
		SecondLevel: random.HardQuestions(5, false),
		ThirdLevel:  random.SimpleQuestions(5, false),
	}

	qSrv := question.New(&cfg, nil, nil)
	qSrv.StartFirstLvl(1)

	controller := New(nil, 0, &cfg, qSrv)

	ctrl := gomock.NewController(t)
	telectx := mocks.NewMockteleCtx(ctrl)

	telectx.EXPECT().EditOrSend(gomock.Any(), gomock.Any()).Return(nil).Do(func(msg string, kb *telebot.ReplyMarkup) {
		expectedMsg := message.SecondLvlMessage

		expectedKb := view.StartSecondLevel()

		assert.Equal(t, expectedMsg, msg)
		assert.Equal(t, expectedKb, kb)
	})

	chat := telebot.Chat{
		ID: 1,
	}

	telectx.EXPECT().Chat().Return(&chat).Times(2)
	telectx.EXPECT().Message().Return(&telebot.Message{})

	err := controller.SendLevelMessage(context.TODO(), telectx)
	require.NoError(t, err)
}

func TestSendLevelMessage_FirstLvl_WithPicture(t *testing.T) {
	cfg := config.Config{
		FirstLevel:  random.SimpleQuestions(1, false),
		SecondLevel: random.HardQuestions(5, false),
		ThirdLevel:  random.SimpleQuestions(5, false),
	}

	qSrv := question.New(&cfg, nil, nil)
	qSrv.StartFirstLvl(1)

	controller := New(nil, 0, &cfg, qSrv)

	ctrl := gomock.NewController(t)
	telectx := mocks.NewMockteleCtx(ctrl)

	telectx.EXPECT().Delete().Return(nil)

	telectx.EXPECT().Send(gomock.Any(), gomock.Any()).Return(nil).Do(func(msg string, kb *telebot.ReplyMarkup) {
		expectedMsg := message.SecondLvlMessage

		expectedKb := view.StartSecondLevel()

		assert.Equal(t, expectedMsg, msg)
		assert.Equal(t, expectedKb, kb)
	})

	chat := telebot.Chat{
		ID: 1,
	}

	telectx.EXPECT().Chat().Return(&chat).Times(2)
	telectx.EXPECT().Message().Return(&telebot.Message{Caption: "test"})

	err := controller.SendLevelMessage(context.TODO(), telectx)
	require.NoError(t, err)
}

func TestSendLevelMessage_SecondLvl(t *testing.T) {
	cfg := config.Config{
		FirstLevel:  random.SimpleQuestions(1, false),
		SecondLevel: random.HardQuestions(5, false),
		ThirdLevel:  random.SimpleQuestions(5, false),
	}

	qSrv := question.New(&cfg, nil, nil)
	qSrv.StartFirstLvl(1)
	qSrv.StartSecondLvl(1)

	controller := New(nil, 0, &cfg, qSrv)

	ctrl := gomock.NewController(t)
	telectx := mocks.NewMockteleCtx(ctrl)

	telectx.EXPECT().EditOrSend(gomock.Any(), gomock.Any()).Return(nil).Do(func(msg string, kb *telebot.ReplyMarkup) {
		expectedMsg := message.ThirdLvlMessage

		expectedKb := view.StartThirdLevel()

		assert.Equal(t, expectedMsg, msg)
		assert.Equal(t, expectedKb, kb)
	})

	chat := telebot.Chat{
		ID: 1,
	}

	telectx.EXPECT().Chat().Return(&chat).Times(2)
	telectx.EXPECT().Message().Return(&telebot.Message{})

	err := controller.SendLevelMessage(context.TODO(), telectx)
	require.NoError(t, err)
}

func TestSendLevelMessage_SecondLvl_WithPicture(t *testing.T) {
	cfg := config.Config{
		FirstLevel:  random.SimpleQuestions(1, false),
		SecondLevel: random.HardQuestions(5, false),
		ThirdLevel:  random.SimpleQuestions(5, false),
	}

	qSrv := question.New(&cfg, nil, nil)
	qSrv.StartFirstLvl(1)
	qSrv.StartSecondLvl(1)

	controller := New(nil, 0, &cfg, qSrv)

	ctrl := gomock.NewController(t)
	telectx := mocks.NewMockteleCtx(ctrl)

	telectx.EXPECT().Delete().Return(nil)

	telectx.EXPECT().Send(gomock.Any(), gomock.Any()).Return(nil).Do(func(msg string, kb *telebot.ReplyMarkup) {
		expectedMsg := message.ThirdLvlMessage

		expectedKb := view.StartThirdLevel()

		assert.Equal(t, expectedMsg, msg)
		assert.Equal(t, expectedKb, kb)
	})

	chat := telebot.Chat{
		ID: 1,
	}

	telectx.EXPECT().Chat().Return(&chat).Times(2)
	telectx.EXPECT().Message().Return(&telebot.Message{Caption: "test"})

	err := controller.SendLevelMessage(context.TODO(), telectx)
	require.NoError(t, err)
}

func TestSendLevelMessage_ThirdLvl(t *testing.T) {
	cfg := config.Config{
		FirstLevel:  random.SimpleQuestions(1, false),
		SecondLevel: random.HardQuestions(5, false),
		ThirdLevel:  random.SimpleQuestions(5, false),
	}

	bot, err := telebot.NewBot(telebot.Settings{
		Token:   "8126298325:AAG8b4ljktyfnUwsFnV1UxOai4Ma1S9eLuw",
		Offline: true,
	})
	require.NoError(t, err)

	ctrl := gomock.NewController(t)
	telectx := mocks.NewMockteleCtx(ctrl)
	db := mocks.NewMockstorage(ctrl)

	qSrv := question.New(&cfg, db, nil)
	qSrv.StartFirstLvl(1)
	qSrv.StartSecondLvl(1)
	qSrv.StartThirdLvl(1)

	controller := New(bot, -1002285490468, &cfg, qSrv)

	telectx.EXPECT().EditOrSend(gomock.Any(), gomock.Any()).Return(nil).Do(func(msg string, kb *telebot.ReplyMarkup) {
		res, err := qSrv.Results(1)
		require.NoError(t, err)

		expectedMsg := fmt.Sprintf(message.Result, res.RigthAnswers[0], len(cfg.FirstLevel),
			res.RigthAnswers[1], len(cfg.SecondLevel),
			res.RigthAnswers[2], len(cfg.ThirdLevel),
			fmt.Sprintf("%.2fs", res.Seconds))

		expectedKb := view.ResultMenu()

		assert.Equal(t, fmt.Sprintf(message.ResultMessage, expectedMsg), msg)
		assert.Equal(t, expectedKb, kb)
	})

	db.EXPECT().SaveResults(gomock.Any(), gomock.Any()).Return(nil)

	chat := telebot.Chat{
		ID:       1,
		Username: "test",
	}

	telectx.EXPECT().Chat().Return(&chat).Times(6)
	telectx.EXPECT().Message().Return(&telebot.Message{})

	err = controller.SendLevelMessage(context.TODO(), telectx)
	require.NoError(t, err)
}

func TestSendLevelMessage_ThirdLvl_WithPicture(t *testing.T) {
	cfg := config.Config{
		FirstLevel:  random.SimpleQuestions(1, false),
		SecondLevel: random.HardQuestions(5, false),
		ThirdLevel:  random.SimpleQuestions(5, false),
	}

	bot, err := telebot.NewBot(telebot.Settings{
		Token:   "8126298325:AAG8b4ljktyfnUwsFnV1UxOai4Ma1S9eLuw",
		Offline: true,
	})
	require.NoError(t, err)

	ctrl := gomock.NewController(t)
	telectx := mocks.NewMockteleCtx(ctrl)
	db := mocks.NewMockstorage(ctrl)

	qSrv := question.New(&cfg, db, nil)
	qSrv.StartFirstLvl(1)
	qSrv.StartSecondLvl(1)
	qSrv.StartThirdLvl(1)

	controller := New(bot, -1002285490468, &cfg, qSrv)

	telectx.EXPECT().Delete().Return(nil)

	telectx.EXPECT().EditOrSend(gomock.Any(), gomock.Any()).Return(nil).Do(func(msg string, kb *telebot.ReplyMarkup) {
		res, err := qSrv.Results(1)
		require.NoError(t, err)

		expectedMsg := fmt.Sprintf(message.Result, res.RigthAnswers[0], len(cfg.FirstLevel),
			res.RigthAnswers[1], len(cfg.SecondLevel),
			res.RigthAnswers[2], len(cfg.ThirdLevel),
			fmt.Sprintf("%.2fs", res.Seconds))

		expectedKb := view.ResultMenu()

		assert.Equal(t, fmt.Sprintf(message.ResultMessage, expectedMsg), msg)
		assert.Equal(t, expectedKb, kb)
	})

	chat := telebot.Chat{
		ID:       1,
		Username: "test",
	}

	db.EXPECT().SaveResults(gomock.Any(), gomock.Any()).Return(nil)

	telectx.EXPECT().Chat().Return(&chat).Times(6)
	telectx.EXPECT().Message().Return(&telebot.Message{Caption: "test"})

	err = controller.SendLevelMessage(context.TODO(), telectx)
	require.NoError(t, err)
}
