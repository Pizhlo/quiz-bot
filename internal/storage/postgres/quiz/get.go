package quiz

import (
	"context"
	"errors"
	"quiz-mod/internal/model"
)

var (
	ErrNoResults = errors.New("no results found by user ID")
)

func (db *quizRepo) AllResults(ctx context.Context, userID int64) ([]model.Result, error) {
	tx, err := db.tx(ctx)
	if err != nil {
		return nil, err
	}

	rows, err := tx.QueryContext(ctx,
		`select tg_id, 
		first_lvl, total_first_lvl, 
		second_lvl, total_second_lvl, 
		third_lvl, total_third_lvl, 
		duration, date	
		from quizzes.results where tg_id = $1`, userID)
	if err != nil {
		return nil, err
	}

	res := []model.Result{}

	for rows.Next() {
		result := model.Result{}
		result.RigthAnswers = make(map[int]int)
		result.TotalAnswers = make(map[int]int)

		rigthAns := make([]int, 3)
		totalAns := make([]int, 3)

		err = rows.Scan(&result.TgID,
			&rigthAns[model.FirstLevel], &totalAns[model.FirstLevel],
			&rigthAns[model.SecondLevel], &totalAns[model.SecondLevel],
			&rigthAns[model.ThirdLevel], &totalAns[model.ThirdLevel],
			&result.Seconds, &result.Date,
		)
		if err != nil {
			return nil, err
		}

		result.RigthAnswers[model.FirstLevel] = rigthAns[model.FirstLevel]
		result.RigthAnswers[model.SecondLevel] = rigthAns[model.SecondLevel]
		result.RigthAnswers[model.ThirdLevel] = rigthAns[model.ThirdLevel]

		result.TotalAnswers[model.FirstLevel] = totalAns[model.FirstLevel]
		result.TotalAnswers[model.SecondLevel] = totalAns[model.SecondLevel]
		result.TotalAnswers[model.ThirdLevel] = totalAns[model.ThirdLevel]

		res = append(res, result)
	}

	if len(res) == 0 {
		return nil, ErrNoResults
	}

	return res, nil
}
