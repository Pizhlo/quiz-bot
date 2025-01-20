package view

import (
	"fmt"
	"quiz-mod/internal/message"
	"quiz-mod/internal/model"

	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

const (
	dateFieldFormat = "02.01.2006 15:04:05"
)

type ResultView struct {
	pages       []string
	currentPage int
}

func NewResult() *ResultView {
	return &ResultView{pages: make([]string, 0), currentPage: 0}
}

func (s *ResultView) Message(results []model.Result) string {
	page := "<b>Твои результаты:\n\n</b>"

	for i, res := range results {
		msg := ""

		// seconds := time.Duration(res.Seconds) * time.Second
		msg += fmt.Sprintf(message.Result,
			res.RigthAnswers[model.FirstLevel],
			res.TotalAnswers[model.FirstLevel],
			res.RigthAnswers[model.SecondLevel],
			res.TotalAnswers[model.SecondLevel],
			res.RigthAnswers[model.ThirdLevel],
			res.TotalAnswers[model.ThirdLevel],
			// res.Duration.String(),
			fmt.Sprintf("%.2fs", res.Seconds),
			// seconds.String(),
		)

		page += fmt.Sprintf("%s\n\nДата: %s", msg, res.Date.Format(dateFieldFormat))

		if i%5 == 0 {
			s.pages = append(s.pages, page)
			page = ""
		}

	}

	return s.pages[0]
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
