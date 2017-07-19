package main

import "fmt"

// API server environment modes.
const (
	DebugMode   string = "debug"   // for development
	ReleaseMode string = "release" // for production
	TestMode    string = "test"    // for testing
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

var apiPort string

// SetPort sets API port.
func SetPort(port string) {
	apiPort = port
}

// Port returns current api port.
func Port() string {
	return apiPort
}

// API version is set at the start of server.
var apiVersion string

// SetAPIVersion sets current API version.
func SetAPIVersion(v string) {
	apiVersion = v
}

// APIVersion returns the current API version.
func APIVersion() string {
	return apiVersion
}

// sameString returns true if a == b, and false otherwise.
func sameString(a string, b string) bool {
	return a == b
}

// stringInSlice returns true if list contains a string, and false otherwise.
func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
