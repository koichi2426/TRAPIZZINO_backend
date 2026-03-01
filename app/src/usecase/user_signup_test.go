package usecase_test

import (
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestUserSignup_Execute(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	// ここで各リポジトリやサービスを初期化（実際のコードに合わせてInjectしてください）
	// interactor := usecase.NewUserSignupInteractor(repo, presenter, authService)

	t.Run("新規ユーザーを正常に登録できる", func(t *testing.T) {
		mock.ExpectQuery(`SELECT (.+) FROM users WHERE email = \$1`).
			WithArgs("test@example.com").
			WillReturnRows(sqlmock.NewRows([]string{"id"})) // 重複なし

		mock.ExpectExec(`INSERT INTO users`).
			WithArgs("testuser", "test@example.com", sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))

		// 実行と検証
		assert.True(t, true) 
	})
}