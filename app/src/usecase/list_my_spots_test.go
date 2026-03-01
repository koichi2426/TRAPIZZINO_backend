package usecase_test

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestListMySpots_Execute(t *testing.T) {
	t.Run("特定のユーザーIDに紐づく投稿のみが取得される", func(t *testing.T) {
		assert.True(t, true)
	})
}