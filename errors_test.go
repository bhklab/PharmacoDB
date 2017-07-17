package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assert.Equal(t, Error{}, NewError())
}

func TestCreateError(t *testing.T) {
	assert.Panics(t, func() { CreateError(1000, 200, "Test error message") }, "No Panic")
	assert.NotPanics(t, func() { CreateError(ErrorTypePrivate, 200, "Test error message") }, "No Panic")
}

func TestLogPrivateError(t *testing.T) {
	assert.Panics(t, func() { LogPrivateError(ErrorTypePublic, nil) }, "No Panic")
	assert.NotPanics(t, func() { LogPrivateError(ErrorTypePrivate, errors.New("New Error")) }, "No Panic")
}
