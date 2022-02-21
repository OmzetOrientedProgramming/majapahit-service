package checkup

import "github.com/pkg/errors"

var (
	ErrPingDBFailed           = errors.New("Failed to ping DB")
	ErrPostgreSQLNotConnected = errors.New("PostgresDB not connected")
)
