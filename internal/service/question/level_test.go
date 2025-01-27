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

func TestCurrentLevel(t *testing.T) {
	srv := &Question{
		users: make(map[int64]userState),
	}

	type test struct {
		name   string
		state  userState
		result int
	}

	tests := []test{
		{
			name: "first lvl",
			state: userState{
				level: model.FirstLevel,
			},
			result: model.FirstLevel,
		},
		{
			name: "second lvl",
			state: userState{
				level: model.SecondLevel,
			},
			result: model.SecondLevel,
		},
		{
			name: "third lvl",
			state: userState{
				level: model.ThirdLevel,
			},
			result: model.ThirdLevel,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv.users[1] = tt.state

			result, err := srv.CurrentLevel(1)
			require.NoError(t, err)

			assert.Equal(t, tt.result, result)
		})
	}
}

func TestStartFirstLvl(t *testing.T) {
	srv := &Question{
		users:      make(map[int64]userState),
		firstLevel: random.SimpleQuestions(2),
	}

	wayback := time.Date(2024, time.May, 19, 1, 2, 3, 4, time.UTC)

	patch := monkey.Patch(time.Now, func() time.Time { return wayback })
	defer patch.Unpatch()

	result := userState{
		level:        model.FirstLevel,
		rigthAnswers: 0,
		maxQuestions: 2,
		question:     0,
		startTime:    time.Now(),
		result: model.Result{
			RigthAnswers: map[int]int{
				model.FirstLevel:  0,
				model.SecondLevel: 0,
				model.ThirdLevel:  0,
			},
			TotalAnswers: map[int]int{
				model.FirstLevel:  2,
				model.SecondLevel: 0,
				model.ThirdLevel:  0,
			},
		},
	}

	srv.StartFirstLvl(1)

	state, err := srv.stateByUser(1)
	require.NoError(t, err)

	assert.Equal(t, result, state)
}

func TestStartSecondLvl(t *testing.T) {
	srv := &Question{
		users:       make(map[int64]userState),
		secondLevel: random.HardQuestions(2),
	}

	wayback := time.Date(2024, time.May, 19, 1, 2, 3, 4, time.UTC)

	patch := monkey.Patch(time.Now, func() time.Time { return wayback })
	defer patch.Unpatch()

	result := userState{
		level:        model.SecondLevel,
		rigthAnswers: 0,
		maxQuestions: 2,
		question:     0,
		startTime:    time.Now(),
		result: model.Result{
			RigthAnswers: map[int]int{
				model.FirstLevel:  0,
				model.SecondLevel: 0,
				model.ThirdLevel:  0,
			},
			TotalAnswers: map[int]int{
				model.FirstLevel:  2,
				model.SecondLevel: 2,
				model.ThirdLevel:  0,
			},
		},
	}

	srv.users[1] = userState{
		level:        model.FirstLevel,
		rigthAnswers: 0,
		maxQuestions: 2,
		question:     0,
		startTime:    time.Now(),
		result: model.Result{
			RigthAnswers: map[int]int{
				model.FirstLevel:  0,
				model.SecondLevel: 0,
				model.ThirdLevel:  0,
			},
			TotalAnswers: map[int]int{
				model.FirstLevel:  2,
				model.SecondLevel: 0,
				model.ThirdLevel:  0,
			},
		},
	}

	srv.StartSecondLvl(1)

	state, err := srv.stateByUser(1)
	require.NoError(t, err)

	assert.Equal(t, result, state)
}

func TestStartThirdLvl(t *testing.T) {
	srv := &Question{
		users:      make(map[int64]userState),
		thirdLevel: random.SimpleQuestions(3),
	}

	wayback := time.Date(2024, time.May, 19, 1, 2, 3, 4, time.UTC)

	patch := monkey.Patch(time.Now, func() time.Time { return wayback })
	defer patch.Unpatch()

	result := userState{
		level:        model.ThirdLevel,
		rigthAnswers: 0,
		maxQuestions: 3,
		question:     0,
		startTime:    time.Now(),
		result: model.Result{
			RigthAnswers: map[int]int{
				model.FirstLevel:  0,
				model.SecondLevel: 0,
				model.ThirdLevel:  0,
			},
			TotalAnswers: map[int]int{
				model.FirstLevel:  2,
				model.SecondLevel: 2,
				model.ThirdLevel:  3,
			},
		},
	}

	srv.users[1] = userState{
		level:        model.SecondLevel,
		rigthAnswers: 0,
		maxQuestions: 2,
		question:     0,
		startTime:    time.Now(),
		result: model.Result{
			RigthAnswers: map[int]int{
				model.FirstLevel:  0,
				model.SecondLevel: 0,
				model.ThirdLevel:  0,
			},
			TotalAnswers: map[int]int{
				model.FirstLevel:  2,
				model.SecondLevel: 2,
				model.ThirdLevel:  0,
			},
		},
	}

	srv.StartThirdLvl(1)

	state, err := srv.stateByUser(1)
	require.NoError(t, err)

	assert.Equal(t, result, state)
}

func TestLevelResults(t *testing.T) {
	srv := &Question{
		users: make(map[int64]userState),
	}

	state := userState{
		rigthAnswers: 3,
	}

	srv.users[1] = state

	result, err := srv.LevelResults(1)
	require.NoError(t, err)

	assert.Equal(t, state.rigthAnswers, result)
}
