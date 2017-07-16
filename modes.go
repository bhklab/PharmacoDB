package main

import "log"

// API server environment modes.
const (
	DebugMode   string = "debug"   // for development
	ReleaseMode string = "release" // for production
	TestMode    string = "test"    // for testing
)

var apiMode = DebugMode

// SetMode sets server environment mode.
func SetMode(e string) {
	switch e {
	case DebugMode:
		apiMode = DebugMode
	case ReleaseMode:
		apiMode = ReleaseMode
	case TestMode:
		apiMode = TestMode
	default:
		log.Fatal("API mode unknown: " + e)
	}
}

// Mode returns current server environment mode.
func Mode() string {
	return apiMode
}
