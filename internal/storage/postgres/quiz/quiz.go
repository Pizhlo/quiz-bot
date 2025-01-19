package quiz

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type quizRepo struct {
	db        *sql.DB
	currentTx *sql.Tx
}

func New(dbURl string) (*quizRepo, error) {
	db, err := sql.Open("postgres", dbURl)
	if err != nil {
		return nil, fmt.Errorf("connect open a db driver: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("cannot connect to a db: %w", err)
	}

	logrus.Infof("successfully connected db on: %s", dbURl)

	return &quizRepo{db: db}, nil
}

func (db *quizRepo) Close() {
	if err := db.db.Close(); err != nil {
		logrus.Errorf("error on closing user repo: %v", err)
	}
}

func (db *quizRepo) tx(ctx context.Context) (*sql.Tx, error) {
	if db.currentTx != nil {
		return db.currentTx, nil
	}

	tx, err := db.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return nil, err
	}

	db.currentTx = tx

	return tx, nil
}

func (db *quizRepo) commit() error {
	tx := db.currentTx
	db.currentTx = nil
	return tx.Commit()
}

func (db *quizRepo) rollback() error {
	tx := db.currentTx
	db.currentTx = nil
	return tx.Rollback()
}
