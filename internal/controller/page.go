package controller

import (
	"strings"

	"gopkg.in/telebot.v3"
)

// NextPage обрабатывает кнопку переключения на следующую страницу
func (c *Controller) NextPage(telectx telebot.Context) error {
	msg, kb := c.questionSrv.NextPage(telectx.Chat().ID)

	err := telectx.Edit(msg, &telebot.SendOptions{
		ReplyMarkup: kb,
	})

	// если пришла ошибка о том, что сообщение не изменено - игнорируем.
	// такая ошибка происходит, если быть на первой странице и нажать кнопку "первая страница".
	// то же самое происходит и с последней страницей
	if err != nil {
		return checkError(err)
	}

	return nil
}

// PrevPage обрабатывает кнопку переключения на предыдущую страницу
func (c *Controller) PrevPage(telectx telebot.Context) error {
	msg, kb := c.questionSrv.PrevPage(telectx.Chat().ID)

	err := telectx.Edit(msg, &telebot.SendOptions{
		ReplyMarkup: kb,
	})

	// если пришла ошибка о том, что сообщение не изменено - игнорируем.
	// такая ошибка происходит, если быть на первой странице и нажать кнопку "первая страница".
	// то же самое происходит и с последней страницей
	if err != nil {
		return checkError(err)
	}

	return nil
}

// LastPage обрабатывает кнопку переключения на последнюю страницу
func (c *Controller) LastPage(telectx telebot.Context) error {
	msg, kb := c.questionSrv.LastPage(telectx.Chat().ID)

	err := telectx.Edit(msg, &telebot.SendOptions{
		ReplyMarkup: kb,
	})

	// если пришла ошибка о том, что сообщение не изменено - игнорируем.
	// такая ошибка происходит, если быть на первой странице и нажать кнопку "первая страница".
	// то же самое происходит и с последней страницей
	if err != nil {
		return checkError(err)
	}

	return nil
}

// FirstPage обрабатывает кнопку переключения на первую страницу
func (c *Controller) FirstPage(telectx telebot.Context) error {
	msg, kb := c.questionSrv.FirstPage(telectx.Chat().ID)

	err := telectx.Edit(msg, &telebot.SendOptions{
		ReplyMarkup: kb,
	})

	// если пришла ошибка о том, что сообщение не изменено - игнорируем.
	// такая ошибка происходит, если быть на первой странице и нажать кнопку "первая страница".
	// то же самое происходит и с последней страницей
	if err != nil {
		return checkError(err)
	}

	return nil
}

func checkError(err error) error {
	switch t := err.(type) {
	case *telebot.Error:
		if strings.Contains(t.Description, "message is not modified") {
			return nil
		}
	default:
		return err
	}

	return err
}
