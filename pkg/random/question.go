package random

import "quiz-bot/internal/model"

func SimpleQuestions(n int) []model.SimpleQuestion {
	res := []model.SimpleQuestion{}

	for i := 0; i < n; i++ {
		res = append(res, SimpleQuestion())
	}

	return res
}

func SimpleQuestion() model.SimpleQuestion {
	return model.SimpleQuestion{
		Question:    *Question(),
		RigthAnswer: String(10),
	}
}

func HardQuestions(n int) []*model.HardQuestion {
	res := []*model.HardQuestion{}

	for i := 0; i < n; i++ {
		res = append(res, HardQuestion())
	}

	return res
}

func HardQuestion() *model.HardQuestion {
	return &model.HardQuestion{
		Question:     *Question(),
		RigthAnswers: Strings(6),
	}
}

func Questions(n int) []model.Question {
	res := []model.Question{}

	for i := 0; i < n; i++ {
		res = append(res, *Question())
	}

	return res
}

func Question() *model.Question {
	return &model.Question{
		Text:    String(5),
		Answers: Strings(6),
		Picture: String(10),
	}
}
