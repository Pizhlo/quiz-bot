package question

import "gopkg.in/telebot.v3"

// NextPage обрабатывает кнопку переключения на следующую страницу
func (s *Question) NextPage(userID int64) (string, *telebot.ReplyMarkup) {
	return s.views[userID].Next(), s.views[userID].Keyboard()
}

// PrevPage обрабатывает кнопку переключения на предыдущую страницу
func (s *Question) PrevPage(userID int64) (string, *telebot.ReplyMarkup) {
	return s.views[userID].Previous(), s.views[userID].Keyboard()
}

// LastPage обрабатывает кнопку переключения на последнюю страницу
func (s *Question) LastPage(userID int64) (string, *telebot.ReplyMarkup) {
	return s.views[userID].Last(), s.views[userID].Keyboard()
}

// FirstPage обрабатывает кнопку переключения на первую страницу
func (s *Question) FirstPage(userID int64) (string, *telebot.ReplyMarkup) {
	return s.views[userID].First(), s.views[userID].Keyboard()
}
