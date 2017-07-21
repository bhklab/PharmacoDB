package api

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

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
	// Set gin mode as well.
	gin.SetMode(mode)
}

// GetModeFromEnv gets mode from environment variable MODE.
// Panics if mode environment variable is not available.
func GetModeFromEnv() string {
	m := os.Getenv("MODE")
	if m == "" {
		panic("MODE environment variable does not exist.")
	}
	return m
}

// Mode returns current server environment mode.
func Mode() string {
	return apiMode
}

// SetPort sets API port.
func SetPort(port string) {
	apiPort = port
}

// GetPortFromEnv gets port from environment variable PORT.
// Panics if port environment variable is not available.
func GetPortFromEnv() string {
	p := os.Getenv("PORT")
	if p == "" {
		panic("PORT environment variable does not exist.")
	}
	return p
}

// Port returns current API port.
func Port() string {
	return apiPort
}

// SetVersion sets current API version.
func SetVersion(v string) {
	apiVersion = v
}

// GetVersionFromEnv gets version from environment variable VERSION.
// Panics if version environment variable is not available.
func GetVersionFromEnv() string {
	v := os.Getenv("VERSION")
	if v == "" {
		panic("VERSION environment variable does not exist.")
	}
	return v
}

// Version returns the current API version.
func Version() string {
	return apiVersion
}
