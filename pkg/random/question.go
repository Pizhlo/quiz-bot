package random

import "quiz-bot/internal/model"

func SimpleQuestions(n int, withPicture bool) []model.SimpleQuestion {
	res := []model.SimpleQuestion{}

	for i := 0; i < n; i++ {
		res = append(res, SimpleQuestion(withPicture))
	}

	return res
}

func SimpleQuestion(withPicture bool) model.SimpleQuestion {
	return model.SimpleQuestion{
		Question:    *Question(withPicture),
		RigthAnswer: String(10),
	}
}

func HardQuestions(n int, withPicture bool) []*model.HardQuestion {
	res := []*model.HardQuestion{}

	for i := 0; i < n; i++ {
		res = append(res, HardQuestion(withPicture))
	}

	return res
}

func HardQuestion(withPicture bool) *model.HardQuestion {
	q := Question(withPicture)
	return &model.HardQuestion{
		Question:     *q,
		RigthAnswers: q.Answers[:2],
	}
}

func Questions(n int, withPicture bool) []model.Question {
	res := []model.Question{}

	for i := 0; i < n; i++ {
		res = append(res, *Question(withPicture))
	}

	return res
}

func Question(withPicture bool) *model.Question {
	q := &model.Question{
		Text:    String(5),
		Answers: Strings(6),
	}

	if withPicture {
		q.Picture = String(10)
	}

	return q
}
