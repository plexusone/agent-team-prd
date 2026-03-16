package main

import (
	"os"

	"github.com/plexusone/agent-team-prd/cmd/prdtool/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
