package util

import "strings"

// ErrorUnwrap function will unwrap all the error and return the list of error
func ErrorUnwrap(err error) ([]string, string) {
	errString := strings.Split(err.Error(), ":")
	errList, errMessage := errString[0], errString[1]

	return strings.Split(errList, ","), strings.TrimSpace(errMessage)
}
