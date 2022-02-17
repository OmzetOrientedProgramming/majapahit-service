package postgres

import (
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestInitSuccess(t *testing.T) {
	_ = godotenv.Load("../../.env")
	db := Init()
	assert.NotNil(t, db)
}

func TestInitFailed(t *testing.T) {
	t.Setenv("DB_USERNAME", "wrongDBname")
	db := Init()
	assert.Nil(t, db)
}
