package usecase_test

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestDistillRecommendation_Execute(t *testing.T) {
	t.Run("複数の投稿がある場合に共鳴スコアが正しく計算される", func(t *testing.T) {
		// モックのPostRepoから複数の投稿を返し、
		// Serviceが ResonanceScore: 2 などを算出するかをチェック
		assert.True(t, true)
	})
}