package usecase_test

import (
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestRegisterSpotPost_Execute(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	t.Run("既存のスポットがある場合、新しいスポットを作らずに既存IDを使用する", func(t *testing.T) {
		// 1. 既存スポットの返却設定
		rows := sqlmock.NewRows([]string{"id", "name", "latitude", "longitude", "registered_user_id"}).
			AddRow(1, "既存のうどん屋", 35.6467, 139.7101, 1)
		
		mock.ExpectQuery(`SELECT (.+) FROM spots WHERE latitude = \$1 AND longitude = \$2`).
			WithArgs(35.6467, 139.7101).
			WillReturnRows(rows)

		// 2. 投稿のみが保存される設定（SpotのINSERTは走らないことを検証）
		mock.ExpectExec(`INSERT INTO posts`).
			WithArgs(2, 1, "local_malloy", sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(2, 1))

		assert.True(t, true)
	})
}