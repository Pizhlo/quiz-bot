package model

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInitRigthAnswers_NilMap(t *testing.T) {
	res := &Result{}

	assert.Nil(t, res.RigthAnswers)

	res.InitRigthAnswers()

	val, ok := res.RigthAnswers[FirstLevel]
	assert.True(t, ok)
	assert.Equal(t, 0, val)

	val, ok = res.RigthAnswers[SecondLevel]
	assert.True(t, ok)
	assert.Equal(t, 0, val)

	val, ok = res.RigthAnswers[ThirdLevel]
	assert.True(t, ok)
	assert.Equal(t, 0, val)
}

func TestInitRigthAnswers_NotNilMap(t *testing.T) {
	res := &Result{
		RigthAnswers: make(map[int]int),
	}

	assert.NotNil(t, res.RigthAnswers)

	res.InitRigthAnswers()

	val, ok := res.RigthAnswers[FirstLevel]
	assert.True(t, ok)
	assert.Equal(t, 0, val)

	val, ok = res.RigthAnswers[SecondLevel]
	assert.True(t, ok)
	assert.Equal(t, 0, val)

	val, ok = res.RigthAnswers[ThirdLevel]
	assert.True(t, ok)
	assert.Equal(t, 0, val)
}

func TestInitTotalAnswers_NilMap(t *testing.T) {
	res := &Result{}

	assert.Nil(t, res.TotalAnswers)

	res.InitTotalAnswers()

	val, ok := res.TotalAnswers[FirstLevel]
	assert.True(t, ok)
	assert.Equal(t, 0, val)

	val, ok = res.TotalAnswers[SecondLevel]
	assert.True(t, ok)
	assert.Equal(t, 0, val)

	val, ok = res.TotalAnswers[ThirdLevel]
	assert.True(t, ok)
	assert.Equal(t, 0, val)
}

func TestInitTotalAnswers_NotNilMap(t *testing.T) {
	res := &Result{
		TotalAnswers: make(map[int]int),
	}

	assert.NotNil(t, res.TotalAnswers)

	res.InitTotalAnswers()

	val, ok := res.TotalAnswers[FirstLevel]
	assert.True(t, ok)
	assert.Equal(t, 0, val)

	val, ok = res.TotalAnswers[SecondLevel]
	assert.True(t, ok)
	assert.Equal(t, 0, val)

	val, ok = res.TotalAnswers[ThirdLevel]
	assert.True(t, ok)
	assert.Equal(t, 0, val)
}

func TestSaveAnswers_NotNIlMap(t *testing.T) {
	res := &Result{
		RigthAnswers: make(map[int]int),
	}

	assert.NotNil(t, res.RigthAnswers)

	res.SaveAnswers(FirstLevel, 3)
	res.SaveAnswers(SecondLevel, 4)
	res.SaveAnswers(ThirdLevel, 5)

	val, ok := res.RigthAnswers[FirstLevel]
	assert.True(t, ok)
	assert.Equal(t, 3, val)

	val, ok = res.RigthAnswers[SecondLevel]
	assert.True(t, ok)
	assert.Equal(t, 4, val)

	val, ok = res.RigthAnswers[ThirdLevel]
	assert.True(t, ok)
	assert.Equal(t, 5, val)
}

func TestSaveAnswers_NIlMap(t *testing.T) {
	res := &Result{}

	assert.Nil(t, res.RigthAnswers)

	res.SaveAnswers(FirstLevel, 3)
	res.SaveAnswers(SecondLevel, 4)
	res.SaveAnswers(ThirdLevel, 5)

	val, ok := res.RigthAnswers[FirstLevel]
	assert.True(t, ok)
	assert.Equal(t, 3, val)

	val, ok = res.RigthAnswers[SecondLevel]
	assert.True(t, ok)
	assert.Equal(t, 4, val)

	val, ok = res.RigthAnswers[ThirdLevel]
	assert.True(t, ok)
	assert.Equal(t, 5, val)
}

func TestSaveTotalAnswers_NotNIlMap(t *testing.T) {
	res := &Result{
		TotalAnswers: make(map[int]int),
	}

	assert.NotNil(t, res.TotalAnswers)

	res.SaveTotalAnswers(FirstLevel, 3)
	res.SaveTotalAnswers(SecondLevel, 4)
	res.SaveTotalAnswers(ThirdLevel, 5)

	val, ok := res.TotalAnswers[FirstLevel]
	assert.True(t, ok)
	assert.Equal(t, 3, val)

	val, ok = res.TotalAnswers[SecondLevel]
	assert.True(t, ok)
	assert.Equal(t, 4, val)

	val, ok = res.TotalAnswers[ThirdLevel]
	assert.True(t, ok)
	assert.Equal(t, 5, val)
}

func TestSaveTotalAnswers_NIlMap(t *testing.T) {
	res := &Result{}

	assert.Nil(t, res.TotalAnswers)

	res.SaveTotalAnswers(FirstLevel, 3)
	res.SaveTotalAnswers(SecondLevel, 4)
	res.SaveTotalAnswers(ThirdLevel, 5)

	val, ok := res.TotalAnswers[FirstLevel]
	assert.True(t, ok)
	assert.Equal(t, 3, val)

	val, ok = res.TotalAnswers[SecondLevel]
	assert.True(t, ok)
	assert.Equal(t, 4, val)

	val, ok = res.TotalAnswers[ThirdLevel]
	assert.True(t, ok)
	assert.Equal(t, 5, val)
}

func TestValid(t *testing.T) {
	type test struct {
		name string
		res  *Result
		err  error
	}

	tests := []test{
		{
			name: "tg ID not set",
			res: &Result{
				Duration: 100,
				Seconds:  100,
				RigthAnswers: map[int]int{
					1: 1,
					2: 2,
					3: 3,
				},
				TotalAnswers: map[int]int{
					1: 1,
					2: 2,
					3: 3,
				},
				Date: time.Now(),
			},
			err: fmt.Errorf("tg ID not set"),
		},
		{
			name: "duration is zero",
			res: &Result{
				TgID:    100,
				Seconds: 100,
				RigthAnswers: map[int]int{
					1: 1,
					2: 2,
					3: 3,
				},
				TotalAnswers: map[int]int{
					1: 1,
					2: 2,
					3: 3,
				},
				Date: time.Now(),
			},
			err: fmt.Errorf("duration is zero"),
		},
		{
			name: "date is zero",
			res: &Result{
				TgID:     100,
				Duration: 200,
				Seconds:  100,
				RigthAnswers: map[int]int{
					1: 1,
					2: 2,
					3: 3,
				},
				TotalAnswers: map[int]int{
					1: 1,
					2: 2,
					3: 3,
				},
			},
			err: fmt.Errorf("date is zero"),
		},
		{
			name: "rigth answers for 1 level is not saved",
			res: &Result{
				TgID:     100,
				Duration: 200,
				Seconds:  100,
				RigthAnswers: map[int]int{
					SecondLevel: 2,
					ThirdLevel:  3,
				},
				TotalAnswers: map[int]int{
					FirstLevel:  1,
					SecondLevel: 2,
					ThirdLevel:  3,
				},
				Date: time.Now(),
			},
			err: fmt.Errorf("rigth answers for 1 level is not saved"),
		},
		{
			name: "rigth answers for 2 level is not saved",
			res: &Result{
				TgID:     100,
				Duration: 200,
				Seconds:  100,
				RigthAnswers: map[int]int{
					FirstLevel: 1,
					ThirdLevel: 3,
				},
				TotalAnswers: map[int]int{
					FirstLevel:  1,
					SecondLevel: 2,
					ThirdLevel:  3,
				},
				Date: time.Now(),
			},
			err: fmt.Errorf("rigth answers for 2 level is not saved"),
		},
		{
			name: "rigth answers for 3 level is not saved",
			res: &Result{
				TgID:     100,
				Duration: 200,
				Seconds:  100,
				RigthAnswers: map[int]int{
					FirstLevel:  1,
					SecondLevel: 2,
				},
				TotalAnswers: map[int]int{
					FirstLevel:  1,
					SecondLevel: 2,
					ThirdLevel:  3,
				},
				Date: time.Now(),
			},
			err: fmt.Errorf("rigth answers for 3 level is not saved"),
		},
		{
			name: "total answers for 1 level is not saved",
			res: &Result{
				TgID:     100,
				Duration: 200,
				Seconds:  100,
				RigthAnswers: map[int]int{
					FirstLevel:  1,
					SecondLevel: 2,
					ThirdLevel:  3,
				},
				TotalAnswers: map[int]int{
					SecondLevel: 2,
					ThirdLevel:  3,
				},
				Date: time.Now(),
			},
			err: fmt.Errorf("total answers for 1 level is not saved"),
		},
		{
			name: "total answers for 2 level is not saved",
			res: &Result{
				TgID:     100,
				Duration: 200,
				Seconds:  100,
				RigthAnswers: map[int]int{
					FirstLevel:  1,
					SecondLevel: 2,
					ThirdLevel:  3,
				},
				TotalAnswers: map[int]int{
					FirstLevel: 1,
					ThirdLevel: 3,
				},
				Date: time.Now(),
			},
			err: fmt.Errorf("total answers for 2 level is not saved"),
		},
		{
			name: "total answers for 3 level is not saved",
			res: &Result{
				TgID:     100,
				Duration: 200,
				Seconds:  100,
				RigthAnswers: map[int]int{
					FirstLevel:  1,
					SecondLevel: 2,
					ThirdLevel:  3,
				},
				TotalAnswers: map[int]int{
					FirstLevel:  1,
					SecondLevel: 2,
				},
				Date: time.Now(),
			},
			err: fmt.Errorf("total answers for 3 level is not saved"),
		},
		{
			name: "valid",
			res: &Result{
				TgID:     100,
				Duration: 200,
				Seconds:  100,
				RigthAnswers: map[int]int{
					FirstLevel:  1,
					SecondLevel: 2,
					ThirdLevel:  3,
				},
				TotalAnswers: map[int]int{
					FirstLevel:  1,
					SecondLevel: 2,
					ThirdLevel:  3,
				},
				Date: time.Now(),
			},
			err: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.res.Valid()
			if tt.err != nil {
				assert.EqualError(t, err, tt.err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
