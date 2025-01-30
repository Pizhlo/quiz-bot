package question

import (
	"context"
	"quiz-bot/internal/model"
	"quiz-bot/mocks"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaveLvlResults(t *testing.T) {
	type test struct {
		name   string
		lvl    int
		state  userState
		result userState
	}

	tests := []test{
		{
			name: "first lvl",
			lvl:  model.FirstLevel,
			state: userState{
				level:        model.FirstLevel,
				result:       model.Result{},
				rigthAnswers: 3,
			},
			result: userState{
				level:        model.FirstLevel,
				rigthAnswers: 3,
				result: model.Result{
					RigthAnswers: map[int]int{
						model.FirstLevel:  3,
						model.SecondLevel: 0,
						model.ThirdLevel:  0,
					},
				},
			},
		},
		{
			name: "second lvl",
			lvl:  model.SecondLevel,
			state: userState{
				level:        model.SecondLevel,
				result:       model.Result{},
				rigthAnswers: 2,
			},
			result: userState{
				level:        model.SecondLevel,
				rigthAnswers: 2,
				result: model.Result{
					RigthAnswers: map[int]int{
						model.FirstLevel:  0,
						model.SecondLevel: 2,
						model.ThirdLevel:  0,
					},
				},
			},
		},
		{
			name: "third lvl",
			lvl:  model.ThirdLevel,
			state: userState{
				level:        model.ThirdLevel,
				result:       model.Result{},
				rigthAnswers: 4,
			},
			result: userState{
				level:        model.ThirdLevel,
				rigthAnswers: 4,
				result: model.Result{
					RigthAnswers: map[int]int{
						model.FirstLevel:  0,
						model.SecondLevel: 0,
						model.ThirdLevel:  4,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := &Question{
				users: map[int64]userState{
					1: tt.state,
				},
			}

			err := srv.SaveLvlResults(1)
			require.NoError(t, err)

			actual, err := srv.stateByUser(1)
			require.NoError(t, err)

			assert.Equal(t, tt.result, actual)
		})
	}
}

func TestResults(t *testing.T) {
	state := userState{
		result: model.Result{
			Duration: 300,
			Seconds:  300,
			TgID:     1,
			RigthAnswers: map[int]int{
				1: 3,
				2: 4,
				3: 5,
			},
			TotalAnswers: map[int]int{
				1: 5,
				2: 5,
				3: 5,
			},
			Date: time.Now(),
		},
	}

	srv := Question{
		users: map[int64]userState{
			1: state,
		},
	}

	result, err := srv.Results(1)
	require.NoError(t, err)

	assert.Equal(t, state.result, result)
}

func TestSaveResults(t *testing.T) {
	state := userState{
		result: model.Result{
			Duration: 300,
			Seconds:  300,
			TgID:     1,
			RigthAnswers: map[int]int{
				model.FirstLevel:  3,
				model.SecondLevel: 4,
				model.ThirdLevel:  5,
			},
			TotalAnswers: map[int]int{
				model.FirstLevel:  5,
				model.SecondLevel: 5,
				model.ThirdLevel:  5,
			},
			Date: time.Now(),
		},
	}

	ctrl := gomock.NewController(t)
	storage := mocks.NewMockstorage(ctrl)

	storage.EXPECT().SaveResults(gomock.Any(), gomock.Any()).Return(nil).Do(func(ctx any, res model.Result) {
		assert.Equal(t, state.result, res)
	})

	srv := Question{
		users: map[int64]userState{
			1: state,
		},
		storage: storage,
	}

	err := srv.SaveResults(context.TODO(), 1)
	require.NoError(t, err)
}
