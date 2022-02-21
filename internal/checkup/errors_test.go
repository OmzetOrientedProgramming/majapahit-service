package checkup

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckUpErrorVariable(t *testing.T) {
	assert.Equal(t, ErrPingDBFailed.Error(), "Failed to ping DB")
	assert.Equal(t, ErrPostgreSQLNotConnected.Error(), "PostgresDB not connected")
}
