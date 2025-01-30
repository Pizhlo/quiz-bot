package question

import (
	"context"
	"quiz-bot/internal/model"
	"quiz-bot/internal/view"
	"quiz-bot/mocks"
	"quiz-bot/pkg/random"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsQuestionLast(t *testing.T) {
	type test struct {
		name   string
		state  userState
		result bool
	}

	tests := []test{
		{
			name: "simple case #1",
			state: userState{
				level:        model.FirstLevel,
				question:     5,
				maxQuestions: 10,
			},
			result: false,
		},
		{
			name: "simple case #2",
			state: userState{
				level:        model.FirstLevel,
				question:     5,
				maxQuestions: 6,
			},
			result: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := Question{
				users: map[int64]userState{
					1: tt.state,
				},
			}

			result, err := srv.IsQuestionLast(1)
			require.NoError(t, err)

			assert.Equal(t, tt.result, result)
		})
	}
}

func TestQuestionNum(t *testing.T) {
	state := userState{
		level:        model.SecondLevel,
		question:     5,
		maxQuestions: 10,
		rigthAnswers: 3,
	}

	srv := Question{
		users: map[int64]userState{
			1: state,
		},
	}

	actual, err := srv.QuestionNum(1)
	require.NoError(t, err)

	assert.Equal(t, state.maxQuestions, actual)
}

func TestCurrentQuestion(t *testing.T) {
	type test struct {
		name   string
		state  userState
		result *model.Question
	}

	simpleQuestions := random.SimpleQuestions(5, false)
	hardQuestions := random.HardQuestions(5, false)

	tests := []test{
		{
			name: "simple case #1: first lvl",
			state: userState{
				level:        model.FirstLevel,
				question:     2,
				maxQuestions: 5,
				rigthAnswers: 3,
			},
			result: &simpleQuestions[2].Question,
		},
		{
			name: "simple case #2: second lvl",
			state: userState{
				level:        model.SecondLevel,
				question:     3,
				maxQuestions: 5,
				rigthAnswers: 3,
			},
			result: &hardQuestions[3].Question,
		},
		{
			name: "simple case #3: third lvl",
			state: userState{
				level:        model.ThirdLevel,
				question:     4,
				maxQuestions: 5,
				rigthAnswers: 3,
			},
			result: &simpleQuestions[4].Question,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := Question{
				users: map[int64]userState{
					1: tt.state,
				},
				firstLevel:  simpleQuestions,
				secondLevel: hardQuestions,
				thirdLevel:  simpleQuestions,
			}

			actual, err := srv.CurrentQuestion(1)
			require.NoError(t, err)

			assert.Equal(t, tt.result, actual)
		})
	}
}

func TestAllResults(t *testing.T) {
	ctx := context.TODO()

	ctrl := gomock.NewController(t)

	db := mocks.NewMockstorage(ctrl)

	results := random.Results(2)

	srv := &Question{
		storage: db,
		views:   make(map[int64]*view.ResultView),
	}

	db.EXPECT().AllResults(gomock.Any(), gomock.Any()).Return(results, nil)

	view := view.New()

	expectedMsg := view.Message(results)
	expectedKb := view.Keyboard()

	actualMsg, actualKb, err := srv.AllResults(ctx, 1)
	require.NoError(t, err)

	assert.Equal(t, expectedMsg, actualMsg)
	assert.Equal(t, expectedKb, actualKb)

}
