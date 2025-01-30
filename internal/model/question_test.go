package model

import (
	"fmt"
	"quiz-bot/internal/message"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQuestionText(t *testing.T) {
	q := Question{Text: "test"}

	actual := q.QuestionText(10, 15)
	expected := fmt.Sprintf(message.Question, 10, 15, q.Text)

	assert.Equal(t, expected, actual)
}

func TestSetUserAnswer_SimpleQuestion_NilMap(t *testing.T) {
	q := SimpleQuestion{}

	assert.Nil(t, q.UserAnswer)

	user := int64(1)
	answer := "user_answer"

	q.SetUserAnswer(user, answer)

	assert.NotNil(t, q.UserAnswer)

	val, ok := q.UserAnswer[user]
	assert.True(t, ok)

	assert.Equal(t, answer, val)
}

func TestSetUserAnswer_SimpleQuestion_NotNilMap(t *testing.T) {
	q := SimpleQuestion{
		UserAnswer: make(map[int64]string),
	}

	assert.NotNil(t, q.UserAnswer)

	user := int64(1)
	answer := "user_answer"

	q.SetUserAnswer(user, answer)

	assert.NotNil(t, q.UserAnswer)

	val, ok := q.UserAnswer[user]
	assert.True(t, ok)

	assert.Equal(t, answer, val)
}

func TestReset_SimpleQuestion(t *testing.T) {
	q := SimpleQuestion{
		UserAnswer: make(map[int64]string),
	}

	assert.NotNil(t, q.UserAnswer)

	user1 := int64(1)
	answer1 := "user_answer1"

	user2 := int64(2)
	answer2 := "user_answer2"

	q.SetUserAnswer(user1, answer1)
	q.SetUserAnswer(user2, answer2)

	q.Reset(user1)

	assert.NotNil(t, q.UserAnswer)

	val, ok := q.UserAnswer[user2]
	assert.True(t, ok)

	assert.Equal(t, answer2, val)

	_, ok = q.UserAnswer[user1]
	assert.False(t, ok)
}

func TestValid_SimpleQuestion(t *testing.T) {
	q := SimpleQuestion{
		RigthAnswer: "Test",
	}

	assert.True(t, q.Valid("test"))

	assert.False(t, q.Valid("random string"))
}

func TestAddUserAnswer_NilMap(t *testing.T) {
	q := HardQuestion{}

	assert.Nil(t, q.UserAnswers)

	user := int64(1)
	answer1 := "user answer"

	q.AddUserAnswer(user, answer1)

	assert.NotNil(t, q.UserAnswers)

	val, ok := q.UserAnswers[user]
	assert.True(t, ok)

	assert.Equal(t, []string{answer1}, val)

	answer2 := "user answer 2"

	q.AddUserAnswer(user, answer2)

	assert.NotNil(t, q.UserAnswers)

	val, ok = q.UserAnswers[user]
	assert.True(t, ok)

	assert.Equal(t, []string{answer1, answer2}, val)
}

func TestAddUserAnswer_NotNilMap(t *testing.T) {
	q := HardQuestion{
		UserAnswers: make(map[int64][]string),
	}

	assert.NotNil(t, q.UserAnswers)

	user := int64(1)
	answer1 := "user answer"

	q.AddUserAnswer(user, answer1)

	assert.NotNil(t, q.UserAnswers)

	val, ok := q.UserAnswers[user]
	assert.True(t, ok)

	assert.Equal(t, []string{answer1}, val)

	answer2 := "user answer 2"

	q.AddUserAnswer(user, answer2)

	assert.NotNil(t, q.UserAnswers)

	val, ok = q.UserAnswers[user]
	assert.True(t, ok)

	assert.Equal(t, []string{answer1, answer2}, val)
}

func TestReset_HardQuestion(t *testing.T) {
	q := HardQuestion{
		UserAnswers: make(map[int64][]string),
	}

	assert.NotNil(t, q.UserAnswers)

	user1 := int64(1)
	answer1 := "user_answer1"

	user2 := int64(2)
	answer2 := "user_answer2"

	q.AddUserAnswer(user1, answer1)
	q.AddUserAnswer(user2, answer2)

	q.Reset(user1)

	assert.NotNil(t, q.UserAnswers)

	val, ok := q.UserAnswers[user2]
	assert.True(t, ok)

	assert.Equal(t, []string{answer2}, val)

	_, ok = q.UserAnswers[user1]
	assert.False(t, ok)
}

func TestValid_HardQuestion(t *testing.T) {
	q := HardQuestion{
		UserAnswers:  make(map[int64][]string),
		RigthAnswers: []string{"answer 1", "answer 2"},
	}

	userAnswers := []string{"answer 1", "answer 2"}

	q.UserAnswers[1] = userAnswers

	assert.True(t, q.Valid(1))

	q.UserAnswers[1] = []string{"test1", "test2"}

	assert.False(t, q.Valid(1))
}
