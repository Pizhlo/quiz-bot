package controller

import (
	"fmt"
	"quiz-mod/internal/view"
	"strings"

	"gopkg.in/telebot.v3"
)

func (c *Controller) SimpleAnswer(telectx telebot.Context) error {
	// получаем правильный ответ, чтобы показать пользователю
	rigthAnswer, err := c.questionSrv.RigthAnswer(telectx.Chat().ID)
	if err != nil {
		return err
	}

	// отправляем ответ в обработку
	err = c.questionSrv.SetAnswer(telectx.Chat().ID, telectx.Data())
	if err != nil {
		return err
	}

	text, err := c.questionSrv.Message(telectx.Chat().ID)
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("%s\n\nТвой ответ: %s\nПравильный ответ: %+v", text, telectx.Data(), rigthAnswer[0])

	return telectx.EditOrSend(msg, view.Next())
}

func (c *Controller) Answer(telectx telebot.Context) error {
	// сохраняем ответ в список
	err := c.questionSrv.AddAnswer(telectx.Chat().ID, telectx.Data())
	if err != nil {
		return err
	}

	msg, err := c.questionSrv.Message(telectx.Chat().ID)
	if err != nil {
		return err
	}

	question, err := c.questionSrv.CurrentQuestion(telectx.Chat().ID)
	if err != nil {
		return err
	}

	userAnswers, err := c.questionSrv.UserAnswers(telectx.Chat().ID)
	if err != nil {
		return err
	}

	answers := []string{}

	// нужно обозначить, какие варианты пользователь уже выбрал
	for _, answer := range question.Answers {
		for _, userAns := range userAnswers {
			if answer == userAns {
				answer = fmt.Sprintf("✅%s", answer)
			}
		}

		answers = append(answers, answer)
	}

	menu := view.Answers(answers)

	return telectx.EditOrSend(msg, menu)
}

// SendAnswer обрабатывает кнопку "отправить" второго уровня при множественном выборе
func (c *Controller) SendAnswer(telectx telebot.Context) error {
	// получаем правильный ответ, чтобы показать пользователю
	rigthAnswers, err := c.questionSrv.RigthAnswer(telectx.Chat().ID)
	if err != nil {
		return err
	}

	text, err := c.questionSrv.Message(telectx.Chat().ID)
	if err != nil {
		return err
	}

	answers := ""

	for _, ans := range rigthAnswers {
		answers += fmt.Sprintf("%s, ", ans)
	}

	answers = strings.Trim(answers, ", ")

	userAnswers, err := c.questionSrv.UserAnswers(telectx.Chat().ID)
	if err != nil {
		return err
	}

	userAnsString := ""

	for _, uAns := range userAnswers {
		userAnsString += fmt.Sprintf("%s, ", uAns)
	}

	userAnsString = strings.Trim(userAnsString, ", ")

	err = c.questionSrv.SaveAnswers(telectx.Chat().ID)
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("%s\n\nТвой ответ: %s\nПравильный ответ: %+v", text, userAnsString, answers)

	return telectx.EditOrSend(msg, view.Next())
}

func (c *Controller) OnText(telectx telebot.Context) error {
	lvl, _ := c.questionSrv.CurrentLevel(telectx.Chat().ID) // не надо проверять ошибку - если бот не знает пользователя, не надо реагировать

	switch lvl {
	case 2:
		rigthAnswers, err := c.questionSrv.RigthAnswer(telectx.Chat().ID)
		if err != nil {
			return err
		}

		err = c.questionSrv.SetAnswer(telectx.Chat().ID, telectx.Text())
		if err != nil {
			return err
		}

		text, err := c.questionSrv.Message(telectx.Chat().ID)
		if err != nil {
			return err
		}

		msg := fmt.Sprintf("%s\n\nТвой ответ: %s\nПравильный ответ: %+v", text, telectx.Text(), rigthAnswers[0])

		return telectx.EditOrSend(msg, view.Next())
	default:
		return nil
	}
}
