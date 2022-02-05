package program_test

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/orlangure/gocovsh/internal/program"
	"github.com/stretchr/testify/require"
)

func TestVersion(t *testing.T) {
	version := "1.2.3"
	commit := "abcdef"
	date := time.Now().Format(time.RFC3339)
	buf := bytes.NewBuffer(nil)

	flagSet := flag.NewFlagSet("test", flag.ContinueOnError)
	p := program.New(
		program.WithBuildInfo(version, commit, date),
		program.WithOutput(buf),
		program.WithFlagSet(flagSet, []string{"--version"}),
	)

	require.NoError(t, p.Run())

	expectedVersion := fmt.Sprintf("Version: %s\nCommit: %s\nDate: %s\n", version, commit, date)
	require.Equal(t, expectedVersion, buf.String())
}

func TestLogger(t *testing.T) {
	buf := bytes.NewBuffer(nil)

	f, err := os.CreateTemp("", "test-logger")
	require.NoError(t, err)

	t.Cleanup(func() {
		require.NoError(t, f.Close())
		require.NoError(t, os.Remove(f.Name()))
	})

	flagSet := flag.NewFlagSet("test", flag.ContinueOnError)
	p := program.New(
		program.WithOutput(buf),
		program.WithLogFile(f.Name()),
		program.WithFlagSet(flagSet, nil),
	)

	// in tests, the program fails
	require.Error(t, p.Run())

	logs, err := os.ReadFile(f.Name())
	require.NoError(t, err)
	require.Contains(t, string(logs), "logging to")
}
