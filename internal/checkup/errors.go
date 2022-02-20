package checkup

import "github.com/pkg/errors"

var (
	ErrPingDBFailed   = errors.New("Failed to ping DB")
	ErrDBNotConnected = errors.New("PostgresDB not connected")
)
