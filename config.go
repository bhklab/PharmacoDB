package main

import "fmt"

// API server environment modes.
const (
	DebugMode   string = "debug"   // for development
	ReleaseMode string = "release" // for production
	TestMode    string = "test"    // for testing
)

var (
	apiPort    string
	apiVersion string
)

var apiMode = ReleaseMode

// SetMode sets server environment mode.
func SetMode(mode string) {
	switch mode {
	case DebugMode:
		apiMode = DebugMode
	case ReleaseMode:
		apiMode = ReleaseMode
	case TestMode:
		apiMode = TestMode
	default:
		panic(fmt.Errorf("API mode '%s' not recognized", mode))
	}
}

// Mode returns current server environment mode.
func Mode() string {
	return apiMode
}

// SetPort sets API port.
func SetPort(port string) {
	apiPort = port
}

// Port returns current api port.
func Port() string {
	return apiPort
}

// SetAPIVersion sets current API version.
func SetAPIVersion(v string) {
	apiVersion = v
}

// APIVersion returns the current API version.
func APIVersion() string {
	return apiVersion
}
