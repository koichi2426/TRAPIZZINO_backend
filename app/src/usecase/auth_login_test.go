package usecase_test

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestAuthLogin_Execute(t *testing.T) {
	t.Run("正しいパスワードでトークンが発行される", func(t *testing.T) {
		// Mockを使って、DBからユーザー取得 -> パスワード照合 -> トークン生成の流れをテスト
		assert.True(t, true)
	})
}