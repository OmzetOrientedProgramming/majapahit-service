package util

import (
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestErrorUnwrap(t *testing.T) {
	expectation := []string{
		"error 1", "error 2",
	}
	errMessageExpectation := errors.New("test error")

	err := errors.Wrap(errMessageExpectation, strings.Join(expectation, ","))

	errorList, message := ErrorUnwrap(err)

	assert.Equal(t, expectation, errorList)
	assert.Equal(t, errMessageExpectation.Error(), message)
}
