package question

import (
	"quiz-bot/internal/model"
	"quiz-bot/pkg/random"
	"testing"

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
