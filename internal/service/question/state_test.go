package question

import (
	"quiz-bot/internal/model"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaveState(t *testing.T) {
	srv := Question{
		users: make(map[int64]userState),
	}

	assert.Empty(t, srv.users)

	state := userState{
		level:        model.SecondLevel,
		question:     5,
		maxQuestions: 10,
		rigthAnswers: 3,
	}

	srv.saveState(1, state)

	actual, err := srv.stateByUser(1)
	require.NoError(t, err)

	assert.Equal(t, state, actual)
}

func TestStateByUser(t *testing.T) {
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

	actual, err := srv.stateByUser(1)
	require.NoError(t, err)

	assert.Equal(t, state, actual)
}

func TestRigthAnswers(t *testing.T) {
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

	actual, err := srv.RigthAnswers(1)
	require.NoError(t, err)

	assert.Equal(t, state.rigthAnswers, actual)
}
