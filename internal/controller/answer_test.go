package controller

import (
	"context"
	"fmt"
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

func TestSimpleAnswer_WithoutPicture(t *testing.T) {
	cfg := config.Config{
		FirstLevel:  random.SimpleQuestions(5, false),
		SecondLevel: random.HardQuestions(5, false),
		ThirdLevel:  random.SimpleQuestions(5, false),
	}
	qSrv := question.New(&cfg, nil, nil, "")
	qSrv.StartFirstLvl(1)

	controller := New(nil, 0, &cfg, qSrv)

	ctrl := gomock.NewController(t)
	telectx := mocks.NewMockteleCtx(ctrl)

	data := "test"

	telectx.EXPECT().EditOrSend(gomock.Any(), gomock.Any()).Return(nil).Do(func(msg string, kb *telebot.ReplyMarkup) {
		expectedMsg, err := qSrv.Message(1)
		require.NoError(t, err)

		expectedKb := view.Next()

		assert.Equal(t, fmt.Sprintf("%s\n\nТвой ответ: %s\nПравильный ответ: %+v", expectedMsg, data, cfg.FirstLevel[0].RigthAnswer), msg)
		assert.Equal(t, expectedKb, kb)
	})

	chat := telebot.Chat{
		ID: 1,
	}

	telectx.EXPECT().Chat().Return(&chat).Times(3)
	telectx.EXPECT().Data().Return(data).Times(2)
	telectx.EXPECT().Message().Return(&telebot.Message{})

	err := controller.SimpleAnswer(telectx)
	require.NoError(t, err)
}

func TestSimpleAnswer_WithPicture(t *testing.T) {
	cfg := config.Config{
		FirstLevel:  random.SimpleQuestions(5, false),
		SecondLevel: random.HardQuestions(5, false),
		ThirdLevel:  random.SimpleQuestions(5, false),
	}
	qSrv := question.New(&cfg, nil, nil, "")
	qSrv.StartFirstLvl(1)

	controller := New(nil, 0, &cfg, qSrv)

	ctrl := gomock.NewController(t)
	telectx := mocks.NewMockteleCtx(ctrl)

	data := "test"

	telectx.EXPECT().EditCaption(gomock.Any(), gomock.Any()).Return(nil).Do(func(caption string, sendOpts *telebot.SendOptions) {
		expectedMsg, err := qSrv.Message(1)
		require.NoError(t, err)

		expectedSendOpts := &telebot.SendOptions{
			ReplyMarkup: view.Next(),
			ParseMode:   htmlParseMode,
		}

		assert.Equal(t, fmt.Sprintf("%s\n\nТвой ответ: %s\nПравильный ответ: %+v", expectedMsg, data, cfg.FirstLevel[0].RigthAnswer), caption)
		assert.Equal(t, expectedSendOpts, sendOpts)
	})

	chat := telebot.Chat{
		ID: 1,
	}

	telectx.EXPECT().Chat().Return(&chat).Times(3)
	telectx.EXPECT().Data().Return(data).Times(2)
	telectx.EXPECT().Message().Return(&telebot.Message{Caption: "test"})

	err := controller.SimpleAnswer(telectx)
	require.NoError(t, err)
}

func TestMultipleAnswer_NoAnswersSelected(t *testing.T) {
	cfg := config.Config{
		FirstLevel:  random.SimpleQuestions(5, false),
		SecondLevel: random.HardQuestions(5, false),
		ThirdLevel:  random.SimpleQuestions(5, false),
	}

	qSrv := question.New(&cfg, nil, nil, "")
	qSrv.StartFirstLvl(1)
	qSrv.StartSecondLvl(1)

	controller := New(nil, 0, &cfg, qSrv)

	ctrl := gomock.NewController(t)
	telectx := mocks.NewMockteleCtx(ctrl)

	chat := telebot.Chat{
		ID: 1,
	}

	data := "test"

	telectx.EXPECT().Chat().Return(&chat).Times(4)
	telectx.EXPECT().Data().Return(data)
	telectx.EXPECT().Message().Return(&telebot.Message{})

	telectx.EXPECT().EditOrSend(gomock.Any(), gomock.Any()).Return(nil).Do(func(msg string, kb *telebot.ReplyMarkup) {
		expectedMsg, err := qSrv.Message(1)
		require.NoError(t, err)

		expectedKb := view.Answers(cfg.SecondLevel[0].Answers)

		assert.Equal(t, expectedMsg, msg)
		assert.Equal(t, expectedKb, kb)
	})

	err := controller.MultipleAnswer(telectx)
	require.NoError(t, err)
}

func TestMultipleAnswer_AnswersSelected(t *testing.T) {
	cfg := config.Config{
		FirstLevel:  random.SimpleQuestions(5, false),
		SecondLevel: random.HardQuestions(5, false),
		ThirdLevel:  random.SimpleQuestions(5, false),
	}

	qSrv := question.New(&cfg, nil, nil, "")
	qSrv.StartFirstLvl(1)
	qSrv.StartSecondLvl(1)

	controller := New(nil, 0, &cfg, qSrv)

	ctrl := gomock.NewController(t)
	telectx := mocks.NewMockteleCtx(ctrl)

	chat := telebot.Chat{
		ID: 1,
	}

	data := cfg.SecondLevel[0].Answers[0]

	telectx.EXPECT().Chat().Return(&chat).Times(4)
	telectx.EXPECT().Data().Return(data)
	telectx.EXPECT().Message().Return(&telebot.Message{})

	cfg.SecondLevel[0].UserAnswers = make(map[int64][]string)
	answers := cfg.SecondLevel[0].Answers
	// берем несколько ответов, как будто их пользователь выбрал
	userAnswers := make([]string, 2)
	n := copy(userAnswers, answers[:2])
	assert.Equal(t, 2, n)

	cfg.SecondLevel[0].UserAnswers[1] = userAnswers

	telectx.EXPECT().EditOrSend(gomock.Any(), gomock.Any()).Return(nil).Do(func(msg string, kb *telebot.ReplyMarkup) {
		expectedMsg, err := qSrv.Message(1)
		require.NoError(t, err)

		answers := []string{}

		questionAnswers := cfg.SecondLevel[0].Answers
		for _, answer := range questionAnswers {
			for _, userAns := range userAnswers {
				if answer == userAns {
					answer = fmt.Sprintf("✅%s", answer)
				}
			}

			answers = append(answers, answer)
		}

		expectedKb := view.Answers(answers)

		assert.Equal(t, expectedMsg, msg)
		assert.Equal(t, expectedKb, kb)
	})

	err := controller.MultipleAnswer(telectx)
	require.NoError(t, err)
}

func TestMultipleAnswer_NoAnswersSelected_WithPicture(t *testing.T) {
	cfg := config.Config{
		FirstLevel:  random.SimpleQuestions(5, false),
		SecondLevel: random.HardQuestions(5, false),
		ThirdLevel:  random.SimpleQuestions(5, false),
	}

	qSrv := question.New(&cfg, nil, nil, "")
	qSrv.StartFirstLvl(1)
	qSrv.StartSecondLvl(1)

	controller := New(nil, 0, &cfg, qSrv)

	ctrl := gomock.NewController(t)
	telectx := mocks.NewMockteleCtx(ctrl)

	chat := telebot.Chat{
		ID: 1,
	}

	data := "test"

	telectx.EXPECT().Chat().Return(&chat).Times(4)
	telectx.EXPECT().Data().Return(data)
	telectx.EXPECT().Message().Return(&telebot.Message{Caption: "test"})

	telectx.EXPECT().EditCaption(gomock.Any(), gomock.Any()).Return(nil).Do(func(msg string, sendOpts *telebot.SendOptions) {
		expectedMsg, err := qSrv.Message(1)
		require.NoError(t, err)

		menu := view.Answers(cfg.SecondLevel[0].Answers)
		expectedSendOpts := &telebot.SendOptions{
			ReplyMarkup: menu,
			ParseMode:   htmlParseMode,
		}

		assert.Equal(t, expectedMsg, msg)
		assert.Equal(t, expectedSendOpts, sendOpts)
	})

	err := controller.MultipleAnswer(telectx)
	require.NoError(t, err)
}

func TestMultipleAnswer_AnswersSelected_WithPicture(t *testing.T) {
	cfg := config.Config{
		FirstLevel:  random.SimpleQuestions(5, false),
		SecondLevel: random.HardQuestions(5, false),
		ThirdLevel:  random.SimpleQuestions(5, false),
	}

	qSrv := question.New(&cfg, nil, nil, "")
	qSrv.StartFirstLvl(1)
	qSrv.StartSecondLvl(1)

	controller := New(nil, 0, &cfg, qSrv)

	ctrl := gomock.NewController(t)
	telectx := mocks.NewMockteleCtx(ctrl)

	chat := telebot.Chat{
		ID: 1,
	}

	data := cfg.SecondLevel[0].Answers[0]

	telectx.EXPECT().Chat().Return(&chat).Times(4)
	telectx.EXPECT().Data().Return(data)
	telectx.EXPECT().Message().Return(&telebot.Message{Caption: "test"})

	cfg.SecondLevel[0].UserAnswers = make(map[int64][]string)
	answers := cfg.SecondLevel[0].Answers
	// берем несколько ответов, как будто их пользователь выбрал
	userAnswers := make([]string, 2)
	n := copy(userAnswers, answers[:2])
	assert.Equal(t, 2, n)

	cfg.SecondLevel[0].UserAnswers[1] = userAnswers

	telectx.EXPECT().EditCaption(gomock.Any(), gomock.Any()).Return(nil).Do(func(msg string, sendOpts *telebot.SendOptions) {
		expectedMsg, err := qSrv.Message(1)
		require.NoError(t, err)

		answers := []string{}

		questionAnswers := cfg.SecondLevel[0].Answers
		for _, answer := range questionAnswers {
			for _, userAns := range userAnswers {
				if answer == userAns {
					answer = fmt.Sprintf("✅%s", answer)
				}
			}

			answers = append(answers, answer)
		}

		menu := view.Answers(answers)

		expectedSendOpts := &telebot.SendOptions{
			ReplyMarkup: menu,
			ParseMode:   htmlParseMode,
		}

		assert.Equal(t, expectedMsg, msg)
		assert.Equal(t, expectedSendOpts, sendOpts)
	})

	err := controller.MultipleAnswer(telectx)
	require.NoError(t, err)
}

func TestSendAnswer(t *testing.T) {
	cfg := config.Config{
		FirstLevel:  random.SimpleQuestions(5, false),
		SecondLevel: random.HardQuestions(5, false),
		ThirdLevel:  random.SimpleQuestions(5, false),
	}

	qSrv := question.New(&cfg, nil, nil, "")
	qSrv.StartFirstLvl(1)
	qSrv.StartSecondLvl(1)

	controller := New(nil, 0, &cfg, qSrv)

	ctrl := gomock.NewController(t)
	telectx := mocks.NewMockteleCtx(ctrl)

	chat := telebot.Chat{
		ID: 1,
	}

	cfg.SecondLevel[0].UserAnswers = make(map[int64][]string)
	answers := cfg.SecondLevel[0].Answers
	// берем несколько ответов, как будто их пользователь выбрал
	userAnswers := make([]string, 2)
	n := copy(userAnswers, answers[:2])
	assert.Equal(t, 2, n)

	cfg.SecondLevel[0].UserAnswers[1] = userAnswers

	telectx.EXPECT().Chat().Return(&chat).Times(4)
	telectx.EXPECT().Message().Return(&telebot.Message{})

	telectx.EXPECT().EditOrSend(gomock.Any(), gomock.Any()).Return(nil).Do(func(msg string, kb *telebot.ReplyMarkup) {
		text, err := qSrv.Message(1)
		require.NoError(t, err)

		userAnsString := fmt.Sprintf("%s, %s", userAnswers[0], userAnswers[1])
		answersString := fmt.Sprintf("%s, %s", answers[0], answers[1])
		expectedMsg := fmt.Sprintf("%s\n\nТвой ответ: %s\nПравильный ответ: %+v", text, userAnsString, answersString)

		expectedKb := view.Next()

		assert.Equal(t, expectedMsg, msg)
		assert.Equal(t, expectedKb, kb)
	})

	err := controller.SendAnswer(telectx)
	require.NoError(t, err)
}

func TestSendAnswer_WithPicture(t *testing.T) {
	cfg := config.Config{
		FirstLevel:  random.SimpleQuestions(5, false),
		SecondLevel: random.HardQuestions(5, false),
		ThirdLevel:  random.SimpleQuestions(5, false),
	}

	qSrv := question.New(&cfg, nil, nil, "")
	qSrv.StartFirstLvl(1)
	qSrv.StartSecondLvl(1)

	controller := New(nil, 0, &cfg, qSrv)

	ctrl := gomock.NewController(t)
	telectx := mocks.NewMockteleCtx(ctrl)

	chat := telebot.Chat{
		ID: 1,
	}

	cfg.SecondLevel[0].UserAnswers = make(map[int64][]string)
	answers := cfg.SecondLevel[0].Answers
	// берем несколько ответов, как будто их пользователь выбрал
	userAnswers := make([]string, 2)
	n := copy(userAnswers, answers[:2])
	assert.Equal(t, 2, n)

	cfg.SecondLevel[0].UserAnswers[1] = userAnswers

	telectx.EXPECT().Chat().Return(&chat).Times(4)
	telectx.EXPECT().Message().Return(&telebot.Message{Caption: "Test"})

	telectx.EXPECT().EditCaption(gomock.Any(), gomock.Any()).Return(nil).Do(func(msg string, sendOpts *telebot.SendOptions) {
		text, err := qSrv.Message(1)
		require.NoError(t, err)

		userAnsString := fmt.Sprintf("%s, %s", userAnswers[0], userAnswers[1])
		answersString := fmt.Sprintf("%s, %s", answers[0], answers[1])
		expectedMsg := fmt.Sprintf("%s\n\nТвой ответ: %s\nПравильный ответ: %+v", text, userAnsString, answersString)

		expectedSendOpts := &telebot.SendOptions{
			ReplyMarkup: view.Next(),
			ParseMode:   htmlParseMode,
		}

		assert.Equal(t, expectedMsg, msg)
		assert.Equal(t, expectedSendOpts, sendOpts)
	})

	err := controller.SendAnswer(telectx)
	require.NoError(t, err)
}

func TestOnText(t *testing.T) {
	cfg := config.Config{
		FirstLevel:  random.SimpleQuestions(5, false),
		SecondLevel: random.HardQuestions(5, false),
		ThirdLevel:  random.SimpleQuestions(5, false),
	}

	qSrv := question.New(&cfg, nil, nil, "")
	qSrv.StartFirstLvl(1)
	qSrv.StartSecondLvl(1)
	qSrv.StartThirdLvl(1)

	controller := New(nil, 0, &cfg, qSrv)

	ctrl := gomock.NewController(t)
	telectx := mocks.NewMockteleCtx(ctrl)

	chat := telebot.Chat{
		ID: 1,
	}

	telectx.EXPECT().Chat().Return(&chat).Times(5)
	telectx.EXPECT().Message().Return(&telebot.Message{})
	telectx.EXPECT().Text().Return("test").Times(2)

	telectx.EXPECT().EditOrSend(gomock.Any(), gomock.Any()).Return(nil).Do(func(msg string, kb *telebot.ReplyMarkup) {
		rigthAnswers, err := qSrv.RigthAnswer(1)
		require.NoError(t, err)

		text, err := qSrv.Message(1)
		require.NoError(t, err)

		expectedMsg := fmt.Sprintf("%s\n\nТвой ответ: %s\nПравильный ответ: %+v", text, "test", rigthAnswers[0])

		expectedKb := view.Next()

		assert.Equal(t, expectedMsg, expectedMsg)
		assert.Equal(t, expectedKb, kb)
	})

	err := controller.OnText(context.TODO(), telectx)
	require.NoError(t, err)
}

func TestOnText_WithPicture(t *testing.T) {
	cfg := config.Config{
		FirstLevel:  random.SimpleQuestions(5, false),
		SecondLevel: random.HardQuestions(5, false),
		ThirdLevel:  random.SimpleQuestions(5, false),
	}

	qSrv := question.New(&cfg, nil, nil, "")
	qSrv.StartFirstLvl(1)
	qSrv.StartSecondLvl(1)
	qSrv.StartThirdLvl(1)

	controller := New(nil, 0, &cfg, qSrv)

	ctrl := gomock.NewController(t)
	telectx := mocks.NewMockteleCtx(ctrl)

	chat := telebot.Chat{
		ID: 1,
	}

	telectx.EXPECT().Chat().Return(&chat).Times(5)
	telectx.EXPECT().Message().Return(&telebot.Message{Caption: "test"})
	telectx.EXPECT().Text().Return("test").Times(2)

	telectx.EXPECT().EditCaption(gomock.Any(), gomock.Any()).Return(nil).Do(func(msg string, sendOpts *telebot.SendOptions) {
		rigthAnswers, err := qSrv.RigthAnswer(1)
		require.NoError(t, err)

		text, err := qSrv.Message(1)
		require.NoError(t, err)

		expectedMsg := fmt.Sprintf("%s\n\nТвой ответ: %s\nПравильный ответ: %+v", text, "test", rigthAnswers[0])

		expectedSendOpts := &telebot.SendOptions{
			ReplyMarkup: view.Next(),
			ParseMode:   htmlParseMode,
		}

		assert.Equal(t, expectedMsg, expectedMsg)
		assert.Equal(t, expectedSendOpts, sendOpts)
	})

	err := controller.OnText(context.TODO(), telectx)
	require.NoError(t, err)
}
