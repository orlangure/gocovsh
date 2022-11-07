package main

import (
	"fmt"
	"os"

	"github.com/orlangure/gocovsh/internal/program"
	"github.com/orlangure/gocovsh/internal/styles"
)

var (
	// build information, set by goreleaser.
	version = "dev"
	commit  string
	date    string
)

func main() {
	styles.SetTheme()

	if err := program.New(
		program.WithGoModInfo(),
		program.WithBuildInfo(version, commit, date),
		program.WithLogFile(os.Getenv("GOCOVSH_LOG_FILE")),
	).Run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
