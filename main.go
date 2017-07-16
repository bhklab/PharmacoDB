package main

import "flag"

func main() {
	// Flags
	// -> mode is one of: debug, release, test
	// -> port: default is 8080
	var (
		mode = flag.String("mode", "debug", "environment mode")
		port = flag.String("port", "8080", "server port")
	)

	flag.Parse()

	// Set environment mode.
	SetMode(*mode)

	e := APIConfiguration{Mode: *mode, Port: *port}
	// Start server
	Init(e)
}
