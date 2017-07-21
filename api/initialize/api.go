package main

import (
	"flag"

	"github.com/bhklab/pharmacodb/api"
)

func main() {
	var (
		c api.Context

		modeFromEnv    bool
		portFromEnv    bool
		versionFromEnv bool
	)

	// Flags
	flag.BoolVar(&modeFromEnv, "os-mode", false, "set true if using os environment variables for mode")
	flag.BoolVar(&portFromEnv, "os-port", false, "set true if using os environment variables for port")
	flag.BoolVar(&versionFromEnv, "os-version", false, "set true if using os environment variables for api version")
	flag.StringVar(&c.Mode, "mode", "release", "environment mode")
	flag.StringVar(&c.Port, "port", "4200", "server port")
	flag.StringVar(&c.Version, "version", "v1", "api version")

	flag.Parse()

	if modeFromEnv {
		c.Mode = api.GetModeFromEnv()
	}
	if portFromEnv {
		c.Port = api.GetPortFromEnv()
	}
	if versionFromEnv {
		c.Version = api.GetVersionFromEnv()
	}

	api.Init(&c)
}
