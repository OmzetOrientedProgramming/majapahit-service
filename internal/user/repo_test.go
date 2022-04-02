package user

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestRepo_GetUserIDByLocalID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer mockDB.Close()
		sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
		repo := NewRepo(sqlxDB)

		rows := mock.NewRows([]string{"id"}).AddRow("1")
		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, phone_number, name, status, COALESCE(email, '') as email, COALESCE(firebase_local_id, '') as firebase_local_id, created_at, updated_at FROM users WHERE firebase_local_id=$1")).
			WithArgs("1").
			WillReturnRows(rows)

		userID, err := repo.GetUserIDByLocalID("1")
		assert.Nil(t, err)
		assert.NotEqual(t, 0, userID)
	})

	t.Run("failed internal server error", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer mockDB.Close()
		sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
		repo := NewRepo(sqlxDB)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, phone_number, name, status, COALESCE(email, ''), COALESCE(firebase_local_id, ''), created_at, updated_at FROM users WHERE firebase_local_id=$1")).
			WithArgs("1").
			WillReturnError(ErrInternalServer)

		userID, err := repo.GetUserIDByLocalID("1")
		assert.NotNil(t, err)
		assert.Zero(t, userID)
		assert.Equal(t, ErrInternalServer, errors.Cause(err))
	})

	t.Run("failed sql no rows", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer mockDB.Close()
		sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
		repo := NewRepo(sqlxDB)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, phone_number, name, status, COALESCE(email, '') as email, COALESCE(firebase_local_id, '') as firebase_local_id, created_at, updated_at FROM users WHERE firebase_local_id=$1")).
			WithArgs("1").
			WillReturnError(sql.ErrNoRows)

		userID, err := repo.GetUserIDByLocalID("1")
		assert.NotNil(t, err)
		assert.Zero(t, userID)
		assert.Equal(t, ErrNotFound, errors.Cause(err))
	})
}
