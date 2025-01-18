package controller

import (
	"fmt"
	"quiz-mod/internal/message"
	"quiz-mod/internal/view"

	"gopkg.in/telebot.v3"
)

func (c *Controller) Answer(telectx telebot.Context) error {
	lvl, err := c.questionSrv.CurrentLevel(telectx.Chat().ID)
	if err != nil {
		return err
	}

	// получаем правильный ответ, чтобы показать пользователю
	rigthAnswer, err := c.questionSrv.RigthAnswer(telectx.Chat().ID)
	if err != nil {
		return err
	}

	switch lvl {
	case 0, 2:
		// отправляем ответ в обработку
		err := c.questionSrv.SetAnswer(telectx.Chat().ID, telectx.Data())
		if err != nil {
			return err
		}
	case 1:
		err := c.questionSrv.SetAnswers(telectx.Chat().ID, []string{telectx.Text()})
		if err != nil {
			return err
		}
	}

	text := telectx.Text()

	msg := fmt.Sprintf("%s\n\nПравильный ответ: %+v", text, rigthAnswer)

	return telectx.EditOrSend(msg, view.Next())
}

func (c *Controller) Next(telectx telebot.Context) error {
	last, err := c.questionSrv.IsQuestionLast(telectx.Chat().ID)
	if err != nil {
		return err
	}

	// если вопрос не последний - отправляем следующий вопрос
	if !last {
		c.questionSrv.SetNext(telectx.Chat().ID)

		return c.sendCurrentQuestion(telectx)
	}

	// отправляем сообщение с результатами раунда
	return c.levelResuls(telectx)
}

// sendLevelMessage отправляет сообщение с описанием уровня
func (c *Controller) SendLevelMessage(telectx telebot.Context) error {
	lvl, err := c.questionSrv.CurrentLevel(telectx.Chat().ID)
	if err != nil {
		return err
	}

	switch lvl {
	case 0:
		c.questionSrv.SetNext(telectx.Chat().ID)
		return telectx.EditOrSend(message.SecondLvlMessage, view.StartSecondLevel())
	case 1:
		c.questionSrv.SetNext(telectx.Chat().ID)
		return telectx.EditOrSend(message.SecondLvlMessage, view.StartSecondLevel())
	case 2:
		return c.results(telectx)
	default:
		return fmt.Errorf("unknown lvl: %+v", lvl)
	}
}

func (c *Controller) levelResuls(telectx telebot.Context) error {
	rigthAns, err := c.questionSrv.LevelResults(telectx.Chat().ID)
	if err != nil {
		return err
	}

	questionsNum, err := c.questionSrv.QuestionNum(telectx.Chat().ID)
	if err != nil {
		return err
	}

	msg := fmt.Sprintf(message.LevelEnd, rigthAns, questionsNum)

	return telectx.EditOrSend(msg, view.NewLvl())
}

func (c *Controller) results(telectx telebot.Context) error {
	return telectx.EditOrSend("results", view.BackToMenu())
}
