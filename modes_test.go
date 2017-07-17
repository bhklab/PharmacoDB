package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetMode(t *testing.T) {
	assert.Equal(t, DebugMode, Mode())

	SetMode(ReleaseMode)
	assert.Equal(t, ReleaseMode, Mode())

	SetMode(TestMode)
	assert.Equal(t, TestMode, Mode())

	SetMode(DebugMode)
	assert.Equal(t, DebugMode, Mode())

	assert.Panics(t, func() { SetMode("NotRealMode") }, "The code did not panic")
}
