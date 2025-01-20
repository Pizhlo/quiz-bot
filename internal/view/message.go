package view

import (
	"fmt"
	"quiz-mod/internal/message"
	"quiz-mod/internal/model"

	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

const (
	dateFieldFormat    = "02.01.2006 15:04:05"
	recordCountPerPage = 5
	maxMessageLen      = 4096
)

type ResultView struct {
	pages       []string
	currentPage int
}

func NewResult() *ResultView {
	return &ResultView{pages: make([]string, 0), currentPage: 0}
}

func (s *ResultView) Message(results []model.Result) string {
	s.pages = make([]string, 0)

	res := "<b>Твои результаты:</b>\n\n"

	for i, result := range results {

		txt := s.fillMsg(i+1, result)

		res += fmt.Sprintf("%s\n\n", txt)

		if i%recordCountPerPage == 0 && i > 0 || len(res) == maxMessageLen {
			s.pages = append(s.pages, res)
			res = ""
		}
	}

	if len(s.pages) < 5 && res != "" {
		s.pages = append(s.pages, res)
	}

	s.currentPage = 0

	return s.pages[0]
}

func (S *ResultView) fillMsg(idx int, result model.Result) string {
	seconds := fmt.Sprintf("%.2fс", result.Seconds)

	msg := fmt.Sprintf("<b> %d. Дата: %s</b>\n\n", idx, result.Date.Format(dateFieldFormat))

	msg += fmt.Sprintf(message.Result,
		result.RigthAnswers[model.FirstLevel], result.TotalAnswers[model.FirstLevel],
		result.RigthAnswers[model.SecondLevel], result.TotalAnswers[model.SecondLevel],
		result.RigthAnswers[model.ThirdLevel], result.TotalAnswers[model.ThirdLevel],
		seconds,
		// fmt.Sprintf("\n\n📅Дата: %s", result.Date.Format(dateFieldFormat)),
	)

	return msg
}

var (
	// inline кнопка для переключения на предыдущую страницу
	BtnPrevPgResults = tele.Btn{Text: "<", Unique: "prev_pg_results"}
	// inline кнопка для переключения на следующую страницу
	BtnNextPgResults = tele.Btn{Text: ">", Unique: "next_pg_results"}

	// inline кнопка для переключения на первую страницу
	BtnFirstPgResults = tele.Btn{Text: "<<", Unique: "start_pg_results"}
	// inline кнопка для переключения на последнюю страницу
	BtnLastPgResults = tele.Btn{Text: ">>", Unique: "end_pg_results"}
)

// Keyboard делает клавиатуру для навигации по страницам
func (v *ResultView) Keyboard() *tele.ReplyMarkup {
	menu := &tele.ReplyMarkup{}

	// если страниц 1, клавиатура не нужна
	if v.total() == 1 {
		menu.Inline(
			menu.Row(BtnBackToMenu),
		)
		return menu
	}

	text := fmt.Sprintf("%d / %d", v.current(), v.total())

	btn := menu.Data(text, "")

	menu.Inline(
		menu.Row(BtnFirstPgResults, BtnPrevPgResults, btn, BtnNextPgResults, BtnLastPgResults),
		menu.Row(BtnBackToMenu),
	)

	return menu
}

// Next возвращает следующую страницу сообщений
func (v *ResultView) Next() string {
	logrus.Debugf("ResultView: getting next page. Current: %d\n", v.currentPage)

	if v.currentPage == v.total()-1 {
		logrus.Debugf("ResultView: current page is the last. Setting current page to 0.\n")
		v.currentPage = 0
	} else {
		v.currentPage++
		logrus.Debugf("ResultView: incrementing current page. New value: %d\n", v.currentPage)
	}

	return v.pages[v.currentPage]
}

// Previous возвращает предыдущую страницу сообщений
func (v *ResultView) Previous() string {
	logrus.Debugf("ResultView: getting previous page. Current: %d\n", v.currentPage)

	if v.currentPage == 0 {
		logrus.Debugf("ResultView: previous page is the last. Setting current page to maximum: %d.\n", v.total())
		v.currentPage = v.total() - 1
	} else {
		v.currentPage--
		logrus.Debugf("ResultView: decrementing current page. New value: %d\n", v.currentPage)
	}

	return v.pages[v.currentPage]
}

// Last возвращает последнюю страницу сообщений
func (v *ResultView) Last() string {
	logrus.Debugf("ResultView: getting the last page. Current: %d\n", v.currentPage)

	v.currentPage = v.total() - 1

	return v.pages[v.currentPage]
}

// First возвращает первую страницу сообщений
func (v *ResultView) First() string {
	logrus.Debugf("ResultView: getting the first page. Current: %d\n", v.currentPage)

	v.currentPage = 0

	return v.pages[v.currentPage]
}

// current возвращает номер текущей страницы
func (v *ResultView) current() int {
	return v.currentPage + 1
}

// total возвращает общее количество страниц
func (v *ResultView) total() int {
	return len(v.pages)
}

// SetCurrentToFirst устанавливает текущий номер страницы на 1
func (v *ResultView) SetCurrentToFirst() {
	v.currentPage = 0
}
