package main

import (
	"fmt"
	"os"

	"github.com/orlangure/gocovsh/internal/program"
)

var (
	// build information, set by goreleaser.
	version string
	commit  string
	date    string
)

func main() {
	if err := program.New(
		program.WithBuildInfo(version, commit, date),
		program.WithLogFile(os.Getenv("GOCOVSH_LOG_FILE")),
	).Run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
