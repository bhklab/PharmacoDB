package main

import (
	"flag"
	"os"
)

func main() {
	var (
		osMode  bool
		osPort  bool
		mode    string
		port    string
		version string
	)

	// Flags
	// -> if osMode is set to true, use env vars instead of flags for server mode
	// -> if osPort is set to true, use env vars instead of flags for server port
	// -> mode is one of: debug, release, test
	// -> port: default is 8080
	// -> version: default is v1
	flag.BoolVar(&osMode, "os-mode", false, "set true if using os environment variables for mode")
	flag.BoolVar(&osPort, "os-port", false, "set true if using os environment variables for port")
	flag.StringVar(&mode, "mode", "release", "environment mode")
	flag.StringVar(&port, "port", "8080", "server port")
	flag.StringVar(&version, "version", "v1", "api version")

	flag.Parse()

	if osMode {
		// Use environment variables for mode.
		mode = os.Getenv("MODE")
	}
	if osPort {
		// Use environment variables for port.
		if port = os.Getenv("PORT"); port == "" {
			panic("PORT environment variable needs to be set with appropriate port value")
		}
	}

	// Set environment mode, panic if mode
	// is not recognized.
	SetMode(mode)

	// Set API port.
	SetPort(port)

	// Set API version.
	SetAPIVersion(version)

	// Set DB using environment variables, and
	// panic if fields are missing (not filled).
	SetDB()

	// Start server.
	Init()
}
