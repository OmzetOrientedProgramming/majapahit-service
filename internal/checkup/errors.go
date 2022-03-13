package checkup

import "github.com/pkg/errors"

var (
	// ErrPingDBFailed Will raise if failed to ping database
	ErrPingDBFailed = errors.New("Failed to ping DB")

	// ErrPostgreSQLNotConnected Will raise if failed to connect with PostgreSQL
	ErrPostgreSQLNotConnected = errors.New("PostgresDB not connected")
)
