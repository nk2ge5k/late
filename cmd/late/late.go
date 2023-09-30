package main

import (
	"fmt"
	"os"

	"late"
	"late/cmd/late/cli"
)

var (
	version = "dev"
	commit  = "local"
	date    = ""
)

func main() {
	late.SetBuildInfo(version, commit, date)

	args := os.Args[1:]
	if err := cli.Run(args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
