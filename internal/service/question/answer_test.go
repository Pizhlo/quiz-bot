package question

import (
	"quiz-bot/internal/model"
	"quiz-bot/pkg/random"
	"testing"
	"time"

	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetAnswer(t *testing.T) {
	type test struct {
		name   string
		answer string
		state  userState
		result userState
	}

	simpleQuestions := random.SimpleQuestions(10)
	hardQuestions := random.HardQuestions(10)

	tests := []test{
		{
			name: "first lvl: valid answer",
			state: userState{
				level:        model.FirstLevel,
				question:     5,
				maxQuestions: 10,
				rigthAnswers: 3,
			},
			answer: simpleQuestions[5].RigthAnswer,
			result: userState{
				level:        model.FirstLevel,
				question:     5,
				maxQuestions: 10,
				rigthAnswers: 4,
			},
		},
		{
			name: "first lvl: invalid answer",
			state: userState{
				level:        model.FirstLevel,
				question:     5,
				maxQuestions: 10,
				rigthAnswers: 3,
			},
			answer: random.String(5),
			result: userState{
				level:        model.FirstLevel,
				question:     5,
				maxQuestions: 10,
				rigthAnswers: 3,
			},
		},
		{
			name: "third lvl: valid answer",
			state: userState{
				level:        model.ThirdLevel,
				question:     2,
				maxQuestions: 10,
				rigthAnswers: 3,
			},
			answer: simpleQuestions[2].RigthAnswer,
			result: userState{
				level:        model.ThirdLevel,
				question:     2,
				maxQuestions: 10,
				rigthAnswers: 4,
			},
		},
		{
			name: "third lvl: invalid answer",
			state: userState{
				level:        model.ThirdLevel,
				question:     5,
				maxQuestions: 10,
				rigthAnswers: 3,
			},
			answer: random.String(5),
			result: userState{
				level:        model.ThirdLevel,
				question:     5,
				maxQuestions: 10,
				rigthAnswers: 3,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := &Question{
				users: map[int64]userState{
					1: tt.state,
				},
				firstLevel:  simpleQuestions,
				secondLevel: hardQuestions,
				thirdLevel:  simpleQuestions,
			}

			err := srv.SetAnswer(1, tt.answer)
			require.NoError(t, err)

			state, err := srv.stateByUser(1)
			require.NoError(t, err)

			assert.Equal(t, tt.result, state)
		})
	}
}

func TestSaveAnswers(t *testing.T) {
	type test struct {
		name   string
		answer []string
		state  userState
		result userState
	}

	hardQuestions := random.HardQuestions(10)
	// проверяем на этом вопросе
	hardQuestions[5].UserAnswers = make(map[int64][]string)

	tests := []test{
		{
			name: "second lvl: valid answer",
			state: userState{
				level:        model.SecondLevel,
				question:     5,
				maxQuestions: 10,
				rigthAnswers: 3,
			},
			answer: hardQuestions[5].RigthAnswers,
			result: userState{
				level:        model.SecondLevel,
				question:     5,
				maxQuestions: 10,
				rigthAnswers: 4,
			},
		},
		{
			name: "second lvl: invalid answer",
			state: userState{
				level:        model.SecondLevel,
				question:     5,
				maxQuestions: 10,
				rigthAnswers: 3,
			},
			answer: random.Strings(5),
			result: userState{
				level:        model.SecondLevel,
				question:     5,
				maxQuestions: 10,
				rigthAnswers: 3,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// выставляем сохраненные ответы на вопрос
			hardQuestions[5].UserAnswers[1] = tt.answer

			srv := &Question{
				users: map[int64]userState{
					1: tt.state,
				},
				secondLevel: hardQuestions,
			}

			// сохраняем накопленные ответы
			err := srv.SaveAnswers(1)
			require.NoError(t, err)

			state, err := srv.stateByUser(1)
			require.NoError(t, err)

			assert.Equal(t, tt.result, state)
		})
	}
}

func TestAddAnswer(t *testing.T) {
	type test struct {
		name        string
		answer      string
		state       userState
		userAnswers []string //сохраненные ответы
		result      *model.HardQuestion
	}

	hardQuestions := random.HardQuestions(10)
	hardQuestions[3].UserAnswers = make(map[int64][]string)

	// то, что хотим сохранить
	userAnswer := random.String(5)

	// то, что уже сохранено
	userAnswers := random.Strings(5)

	tests := []test{
		{
			name:   "simple case #1: empty slice of user's answer",
			answer: userAnswer,
			state: userState{
				level:        model.SecondLevel,
				question:     3,
				maxQuestions: 10,
				rigthAnswers: 2,
			},
			result: &model.HardQuestion{
				Question: model.Question{
					Text:    hardQuestions[3].Text,
					Answers: hardQuestions[3].Answers,
					Picture: hardQuestions[3].Picture,
				},
				RigthAnswers: hardQuestions[3].RigthAnswers,
				UserAnswers: map[int64][]string{
					1: {userAnswer},
				},
			},
		},
		{
			name:   "simple case #1: not empty slice of user's answer",
			answer: userAnswer,
			state: userState{
				level:        model.SecondLevel,
				question:     3,
				maxQuestions: 10,
				rigthAnswers: 2,
			},
			userAnswers: userAnswers,
			result: &model.HardQuestion{
				Question: model.Question{
					Text:    hardQuestions[3].Text,
					Answers: hardQuestions[3].Answers,
					Picture: hardQuestions[3].Picture,
				},
				RigthAnswers: hardQuestions[3].RigthAnswers,
				UserAnswers: map[int64][]string{
					1: append(userAnswers, userAnswer),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := &Question{
				secondLevel: hardQuestions,
				users: map[int64]userState{
					1: tt.state,
				},
			}

			hardQuestions[3].UserAnswers[1] = tt.userAnswers

			err := srv.AddAnswer(1, tt.answer)
			require.NoError(t, err)

			actual := hardQuestions[3]

			assert.Equal(t, tt.result, actual)
		})
	}
}

func TestRigthAnswer(t *testing.T) {
	hardQuestions := random.HardQuestions(3)
	simpleQuestions := random.SimpleQuestions(3)

	srv := &Question{
		firstLevel:  simpleQuestions,
		secondLevel: hardQuestions,
		thirdLevel:  simpleQuestions,
		users:       make(map[int64]userState),
	}

	type test struct {
		name   string
		state  userState
		result []string
	}

	tests := []test{
		{
			name: "first lvl",
			state: userState{
				level:    model.FirstLevel,
				question: 1,
			},
			result: []string{simpleQuestions[1].RigthAnswer},
		},
		{
			name: "second lvl",
			state: userState{
				level:    model.SecondLevel,
				question: 1,
			},
			result: hardQuestions[1].RigthAnswers,
		},
		{
			name: "third lvl",
			state: userState{
				level:    model.ThirdLevel,
				question: 2,
			},
			result: []string{simpleQuestions[2].RigthAnswer},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv.users[1] = tt.state

			actual, err := srv.RigthAnswer(1)
			require.NoError(t, err)

			assert.Equal(t, tt.result, actual)
		})
	}
}

func TestUserAnswers(t *testing.T) {
	hardQuestions := random.HardQuestions(3)

	srv := &Question{
		secondLevel: hardQuestions,
		users:       make(map[int64]userState),
	}

	state := userState{
		level:    model.SecondLevel,
		question: 1,
	}

	result := hardQuestions[1].UserAnswers[1]

	srv.users[1] = state

	actual, err := srv.UserAnswers(1)
	require.NoError(t, err)

	assert.Equal(t, result, actual)
}

func TestStopTimer(t *testing.T) {
	wayback := time.Date(2025, 01, 27, 10, 20, 3, 4, time.UTC)

	patch := monkey.Patch(time.Now, func() time.Time { return wayback })
	defer patch.Unpatch()

	srv := &Question{
		users: map[int64]userState{
			1: {
				level:        model.ThirdLevel,
				question:     5,
				maxQuestions: 5,
				rigthAnswers: 3,
				startTime:    time.Date(2025, 01, 27, 10, 15, 3, 4, time.UTC),
			},
		},
	}

	result := userState{
		level:        model.ThirdLevel,
		question:     5,
		maxQuestions: 5,
		rigthAnswers: 3,
		startTime:    time.Date(2025, 01, 27, 10, 15, 3, 4, time.UTC),
		result: model.Result{
			Seconds:  300,
			Duration: 300000000000,
			Date:     time.Now(),
		},
	}

	err := srv.StopTimer(1)
	require.NoError(t, err)

	state, err := srv.stateByUser(1)
	require.NoError(t, err)

	assert.Equal(t, result, state)
}
