package question

import (
	"context"
	"fmt"
	"quiz-bot/internal/config"
	"quiz-bot/internal/model"
	"quiz-bot/internal/view"
)

// сервис, управляющий вопросами
type Question struct {
	firstLevel     []model.SimpleQuestion
	secondLevel    []*model.HardQuestion
	thirdLevel     []model.SimpleQuestion
	successPicFile string // название файла с фотографией для конца викторины

	users map[int64]userState // для хранения состояний пользователей

	views map[int64]*view.ResultView // мапа с вьюхами

	storage storage
	minio   minio
}

//go:generate mockgen -source ./service.go -destination=../../../mocks/minio.go -package=mocks
type minio interface {
	// Get получает файл из minio и сохраняет по пути filepath/objectName
	Get(ctx context.Context, filePath string) error
}

//go:generate mockgen -source ./service.go -destination=../../../mocks/question_srv.go -package=mocks
type storage interface {
	// SaveResults сохраняет результат викторины в БД
	SaveResults(ctx context.Context, res model.Result) error

	// AllResults возвращает все результаты викторин пользователя
	AllResults(ctx context.Context, userID int64) ([]model.Result, error)
}

func New(cfg *config.Config, storage storage, minio minio, successPicFile string) *Question {
	return &Question{
		firstLevel:     cfg.FirstLevel,
		secondLevel:    cfg.SecondLevel,
		thirdLevel:     cfg.ThirdLevel,
		users:          make(map[int64]userState),
		storage:        storage,
		views:          make(map[int64]*view.ResultView),
		minio:          minio,
		successPicFile: successPicFile,
	}
}

func (s *Question) Message(userID int64) (string, error) {
	state, err := s.stateByUser(userID)
	if err != nil {
		return "", err
	}

	curIdx := state.question + 1

	switch state.level {
	case model.FirstLevel:
		question := s.firstLevel[state.question]
		return question.QuestionText(curIdx, state.maxQuestions), nil
	case model.SecondLevel:
		question := s.secondLevel[state.question]
		return question.QuestionText(curIdx, state.maxQuestions), nil
	case model.ThirdLevel:
		question := s.thirdLevel[state.question]
		return question.QuestionText(curIdx, state.maxQuestions), nil
	default:
		return "", fmt.Errorf("unknown level: %+v", state.level)
	}
}

func (s *Question) SetNext(userID int64) error {
	state, err := s.stateByUser(userID)
	if err != nil {
		return err
	}

	if state.isQuestionLast() {
		state.level++
		state.question = 0
	} else {
		state.question++
	}

	s.SaveState(userID, state)

	return nil
}

// Reset стирает все сохраненные данные
func (s *Question) Reset(userID int64) {
	delete(s.users, userID)

	for _, q := range s.firstLevel {
		q.Reset(userID)
	}

	for _, q := range s.secondLevel {
		q.Reset(userID)
	}

	for _, q := range s.thirdLevel {
		q.Reset(userID)
	}
}

// GetPics делает запрос в minio и возвращает путь до сохраненного файла
func (s *Question) GetFile(ctx context.Context, filePath string) error {
	return s.minio.Get(ctx, filePath)
}

// GetFileEndRound достает с сервера картинку, которую нужно приложить в конце раунда.
// Возвращает название файла в формате folder/filename.jpg
func (s *Question) GetFileEndRound(ctx context.Context) (string, error) {
	return s.successPicFile, s.minio.Get(ctx, s.successPicFile)
}
