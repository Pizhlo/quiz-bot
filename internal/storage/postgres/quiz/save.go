package quiz

import (
	"context"
	"quiz-mod/internal/model"
)

func (db *quizRepo) SaveResults(ctx context.Context, res model.Result) error {
	tx, err := db.tx(ctx)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx,
		`insert into quizzes.results (tg_id, first_lvl, second_lvl, third_lvl, duration, date) values 
	($1, $2, $3, $4, $5, $6);`,
		res.TgID, res.RigthAnswers[0], res.RigthAnswers[1], res.RigthAnswers[2], res.Duration, res.Date)
	if err != nil {
		_ = db.rollback()
		return err
	}

	return db.commit()
}
