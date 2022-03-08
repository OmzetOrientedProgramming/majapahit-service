package util

import "strings"

func ErrorUnwrap(err error) ([]string, string) {
	errString := strings.Split(err.Error(), ":")
	errList, errMessage := errString[0], errString[1]

	return strings.Split(errList, ","), strings.TrimSpace(errMessage)
}