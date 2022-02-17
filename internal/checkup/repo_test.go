package checkup

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRepo_GetApplicationCheckUp(t *testing.T) {
	testRepo := NewRepo(nil)

	up, err := testRepo.GetApplicationCheckUp()

	assert.Equal(t, true, up)
	assert.Equal(t, nil, err)
}
