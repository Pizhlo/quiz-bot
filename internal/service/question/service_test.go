package question

import (
	"fmt"
	"quiz-bot/internal/message"
	"quiz-bot/internal/model"
	"quiz-bot/pkg/random"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMessage(t *testing.T) {
	srv := Question{
		firstLevel:  random.SimpleQuestions(3),
		secondLevel: random.HardQuestions(3),
		thirdLevel:  random.SimpleQuestions(4),
		users:       make(map[int64]userState),
	}

	type test struct {
		name   string
		state  userState
		result string
	}

	tests := []test{
		{
			name: "first lvl",
			state: userState{
				level:        model.FirstLevel,
				question:     0,
				maxQuestions: len(srv.firstLevel),
			},
			result: fmt.Sprintf(message.Question, 1, len(srv.firstLevel), srv.firstLevel[0].Text),
		},
		{
			name: "second lvl",
			state: userState{
				level:        model.SecondLevel,
				question:     1,
				maxQuestions: len(srv.secondLevel),
			},
			result: fmt.Sprintf(message.Question, 2, len(srv.secondLevel), srv.secondLevel[1].Text),
		},
		{
			name: "third lvl",
			state: userState{
				level:        model.ThirdLevel,
				question:     2,
				maxQuestions: len(srv.thirdLevel),
			},
			result: fmt.Sprintf(message.Question, 3, len(srv.thirdLevel), srv.thirdLevel[2].Text),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv.users[1] = tt.state

			result, err := srv.Message(1)
			require.NoError(t, err)

			assert.Equal(t, tt.result, result)
		})
	}
}

func TestSetNext(t *testing.T) {
	srv := Question{
		users: make(map[int64]userState),
	}

	type test struct {
		name   string
		state  userState
		result userState
	}

	tests := []test{
		{
			name: "first lvl; question not last",
			state: userState{
				level:        model.FirstLevel,
				question:     2,
				maxQuestions: 4,
			},
			result: userState{
				level:        model.FirstLevel,
				question:     3,
				maxQuestions: 4,
			},
		},
		{
			name: "first lvl; question last",
			state: userState{
				level:        model.FirstLevel,
				question:     3,
				maxQuestions: 4,
			},
			result: userState{
				level:        model.SecondLevel,
				question:     0,
				maxQuestions: 4,
			},
		},
		{
			name: "second lvl; question not last",
			state: userState{
				level:        model.SecondLevel,
				question:     2,
				maxQuestions: 4,
			},
			result: userState{
				level:        model.SecondLevel,
				question:     3,
				maxQuestions: 4,
			},
		},
		{
			name: "second lvl; question last",
			state: userState{
				level:        model.SecondLevel,
				question:     3,
				maxQuestions: 4,
			},
			result: userState{
				level:        model.ThirdLevel,
				question:     0,
				maxQuestions: 4,
			},
		},
		{
			name: "third lvl; question not last",
			state: userState{
				level:        model.ThirdLevel,
				question:     2,
				maxQuestions: 4,
			},
			result: userState{
				level:        model.ThirdLevel,
				question:     3,
				maxQuestions: 4,
			},
		},
		{
			name: "third lvl; question last",
			state: userState{
				level:        model.ThirdLevel,
				question:     3,
				maxQuestions: 4,
			},
			result: userState{
				level:        3,
				question:     0,
				maxQuestions: 4,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv.users[1] = tt.state

			err := srv.SetNext(1)
			require.NoError(t, err)

			actual, err := srv.stateByUser(1)
			require.NoError(t, err)

			assert.Equal(t, tt.result, actual)
		})
	}
}

func TestReset(t *testing.T) {
	simpleQ := random.SimpleQuestions(5)
	hardQ := random.HardQuestions(5)

	srv := Question{
		users: map[int64]userState{
			1: {
				level:        model.SecondLevel,
				question:     2,
				maxQuestions: len(hardQ),
				rigthAnswers: 2,
			},
		},
		firstLevel:  simpleQ,
		secondLevel: hardQ,
		thirdLevel:  simpleQ,
	}

	srv.Reset(1)

	_, ok := srv.users[1]
	assert.False(t, ok)

	for _, q := range srv.firstLevel {
		_, ok := q.UserAnswer[1]
		assert.False(t, ok)
	}

	for _, q := range srv.secondLevel {
		_, ok := q.UserAnswers[1]
		assert.False(t, ok)
	}

	for _, q := range srv.thirdLevel {
		_, ok := q.UserAnswer[1]
		assert.False(t, ok)
	}
}
