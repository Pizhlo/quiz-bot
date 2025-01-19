package model

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

type Result struct {
	TgID         int64
	Duration     time.Duration // чтобы засекать время
	Seconds      float64       // для выгрузки из БД
	RigthAnswers map[int]int   // количество правильных ответов пользователя по раундам
	TotalAnswers map[int]int   // всего вопросов в раунде
	Date         time.Time
}

func (s *Result) InitRigthAnswers(user int64) {
	if s.RigthAnswers == nil {
		logrus.Debugf("Result: making new map of rigth answers")
		s.RigthAnswers = make(map[int]int)

		s.RigthAnswers[FirstLevel] = 0
		s.RigthAnswers[SecondLevel] = 0
		s.RigthAnswers[ThirdLevel] = 0
	}
}

func (s *Result) InitTotalAnswers(user int64) {
	if s.TotalAnswers == nil {
		s.TotalAnswers = make(map[int]int)

		s.TotalAnswers[FirstLevel] = 0
		s.TotalAnswers[SecondLevel] = 0
		s.TotalAnswers[ThirdLevel] = 0
	}
}

// SaveAnswers сохраняет количество правильных вопросов за раунд
func (s *Result) SaveAnswers(userID int64, lvl int, result int) {
	if s.RigthAnswers != nil {
		s.RigthAnswers[lvl] = result
	} else {
		s.InitRigthAnswers(userID)
	}

}

// SaveTotalAnswers сохраняет общее количество вопросов за раунд
func (s *Result) SaveTotalAnswers(userID int64, lvl int, result int) {
	if s.TotalAnswers != nil {
		s.TotalAnswers[lvl] = result
	} else {
		s.InitTotalAnswers(userID)
	}

}

func (s *Result) Valid() error {
	if s.TgID == 0 {
		return fmt.Errorf("tg ID not set")
	}

	if s.Duration == 0 {
		return fmt.Errorf("duration is zero")
	}

	if s.Date.IsZero() {
		return fmt.Errorf("date is zero")
	}

	_, ok := s.RigthAnswers[FirstLevel]
	if !ok {
		return fmt.Errorf("rigth answers for 1 level is not saved")
	}

	_, ok = s.RigthAnswers[SecondLevel]
	if !ok {
		return fmt.Errorf("rigth answers for 2 level is not saved")
	}

	_, ok = s.RigthAnswers[ThirdLevel]
	if !ok {
		return fmt.Errorf("rigth answers for 3 level is not saved")
	}

	_, ok = s.TotalAnswers[FirstLevel]
	if !ok {
		return fmt.Errorf("total answers for 1 level is not saved")
	}

	_, ok = s.TotalAnswers[SecondLevel]
	if !ok {
		return fmt.Errorf("total answers for 2 level is not saved")
	}
	_, ok = s.TotalAnswers[ThirdLevel]
	if !ok {
		return fmt.Errorf("total answers for 3 level is not saved")
	}

	return nil
}
