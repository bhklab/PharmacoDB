package api

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

// Config is an API configuration struct.
type Config struct {
	Mode    string
	Port    string
	Version string
}

// API server environment modes.
const (
	DebugMode   string = "debug"   // for development
	ReleaseMode string = "release" // for production
	TestMode    string = "test"    // for testing
)

// DefaultConfig returns the default API configuration setting.
func DefaultConfig() Config {
	return Config{ReleaseMode, "8080", "1"}
}

// SetMode sets api mode.
func (c *Config) SetMode(mode string) {
	switch mode {
	case DebugMode:
		c.Mode = DebugMode
	case ReleaseMode:
		c.Mode = ReleaseMode
	case TestMode:
		c.Mode = TestMode
	default:
		panic(fmt.Errorf("API mode '%s' not recognized", mode))
	}
	// Set gin mode.
	gin.SetMode(mode)
}

// SetPort sets api port.
func (c *Config) SetPort(port string) {
	c.Port = port
}

// SetVersion sets api version.
func (c *Config) SetVersion(version string) {
	c.Version = version
}

// GetEnvMode gets api mode from environment variable MODE.
func GetEnvMode() string {
	m := os.Getenv("MODE")
	if m == "" {
		panic("MODE environment variable does not exist.")
	}
	return m
}

// GetEnvPort gets api port from environment variable PORT.
func GetEnvPort() string {
	p := os.Getenv("PORT")
	if p == "" {
		panic("PORT environment variable does not exist.")
	}
	return p
}

// GetEnvVersion gets api version from environment variable VERSION.
func GetEnvVersion() string {
	v := os.Getenv("VERSION")
	if v == "" {
		panic("VERSION environment variable does not exist.")
	}
	return v
}
