package checkup

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestRepo_GetApplicationCheckUpSuccess(t *testing.T) {
	// Mock DB
	mockDB, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	// Create repository with mock db
	repoMock := NewRepo(sqlxDB)

	// Create expectation

	mock.ExpectPing().WillDelayFor(1 * time.Second)

	up, err := repoMock.GetApplicationCheckUp()

	assert.Equal(t, true, up)
	assert.Equal(t, nil, err)
}

func TestRepo_GetApplicationCheckUpPingError(t *testing.T) {
	// Mock DB
	mockDB, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	// Create repository with mock db
	repoMock := NewRepo(sqlxDB)

	// Create expectation

	mock.ExpectPing().WillReturnError(ErrPingDBFailed)

	up, err := repoMock.GetApplicationCheckUp()

	assert.Equal(t, false, up)
	assert.Equal(t, ErrPingDBFailed, errors.Cause(err))
}

func TestRepo_GetApplicationCheckUpDBNotConnected(t *testing.T) {
	// Create repository with mock db
	repoMock := NewRepo(nil)

	up, err := repoMock.GetApplicationCheckUp()

	assert.Equal(t, false, up)
	assert.Equal(t, ErrDBNotConnected, errors.Cause(err))
}
