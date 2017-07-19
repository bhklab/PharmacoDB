package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogPrivateError(t *testing.T) {
	assert.Panics(t, func() { LogPrivateError(ErrorTypePublic, nil) }, "No Panic")
	assert.NotPanics(t, func() { LogPrivateError(ErrorTypePrivate, errors.New("New Error")) }, "No Panic")
}

func TestLogPublicError(t *testing.T) {
	assert.Panics(t, func() { LogPublicError(nil, ErrorTypePrivate, 200, "Error Message") }, "No Panic")
}
