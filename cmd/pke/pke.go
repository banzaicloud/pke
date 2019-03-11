package main

import (
	"os"

	"github.com/banzaicloud/pke/cmd/pke/app"
)

var (
	Version      string
	CommitHash   string
	GitTreeState string
	BuildDate    string
)

func main() {
	if err := app.Run(Version, CommitHash, GitTreeState, BuildDate); err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
