package main

import (
	"errors"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type MultiError struct {
	errors []error
}

func (e *MultiError) Error() string {
	if len(e.errors) == 0 {
		return ""
	}
	errToStr := strings.Builder{}
	errToStr.WriteString(strconv.Itoa(len(e.errors)))
	errToStr.WriteString(" errors occured:\n")
	for _, err := range e.errors {
		errToStr.WriteString("\t* " + err.Error())
	}
	errToStr.WriteString("\n")
	return errToStr.String()
}

func Append(err error, errs ...error) *MultiError {
	var multiErr *MultiError
	ok := errors.As(err, &multiErr)
	if !ok {
		multiErr = &MultiError{errs}
		return multiErr
	}
	multiErr.errors = append(multiErr.errors, errs...)
	return multiErr
}

func TestMultiError(t *testing.T) {
	var err error
	err = Append(err, errors.New("error 1"))
	err = Append(err, errors.New("error 2"))

	expectedMessage := "2 errors occured:\n\t* error 1\t* error 2\n"
	assert.EqualError(t, err, expectedMessage)
}
