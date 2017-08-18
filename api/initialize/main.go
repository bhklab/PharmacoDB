package main

import (
	"flag"

	"github.com/bhklab/PharmacoDB/api"
)

func main() {
	var (
		c api.Config

		getModeFromEnv    bool
		getPortFromEnv    bool
		getVersionFromEnv bool
	)

	// Flags
	flag.BoolVar(&getModeFromEnv, "emode", false, "use MODE environment variable")
	flag.BoolVar(&getPortFromEnv, "eport", false, "use PORT environment variable")
	flag.BoolVar(&getVersionFromEnv, "eversion", false, "use VERSION environment variable")
	flag.StringVar(&c.Mode, "mode", "release", "environment mode")
	flag.StringVar(&c.Port, "port", "8080", "server port")
	flag.StringVar(&c.Version, "version", "1", "api version")

	flag.Parse()

	if getModeFromEnv {
		c.Mode = api.GetEnvMode()
	}
	if getPortFromEnv {
		c.Port = api.GetEnvPort()
	}
	if getVersionFromEnv {
		c.Version = api.GetEnvVersion()
	}

	api.Init(&c)
}
