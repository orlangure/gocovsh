package program

import (
	"flag"
	"os"
	"testing"

	"github.com/orlangure/gocovsh/internal/gocovshtest/input"
	"github.com/stretchr/testify/require"
)

const (
	filesList = `
foo.go
bar.go
baz.go
`

	diffStr = `
diff --git a/main.go b/main.go
index e6cd709..adc400e 100644
--- a/main.go
+++ b/main.go
@@ -2,6 +2,7 @@ package main
 
 import (
 	"fmt"
+	_ "foo"
 	"os"
 
 	"github.com/orlangure/gocovsh/internal/program"
`
)

func TestParseInput(t *testing.T) {
	t.Parallel()

	t.Run("files list", func(t *testing.T) {
		flagSet := flag.NewFlagSet("test", flag.ContinueOnError)
		f := input.NewMockFile(filesList, os.ModeNamedPipe)
		p := New(
			WithInput(f),
			WithFlagSet(flagSet, nil),
		)

		err := p.parseInput()
		require.NoError(t, err)
		require.Len(t, p.requestedFiles, 3)
		require.EqualValues(t, []string{"foo.go", "bar.go", "baz.go"}, p.requestedFiles)
	})

	t.Run("diff", func(t *testing.T) {
		flagSet := flag.NewFlagSet("test", flag.ContinueOnError)
		f := input.NewMockFile(diffStr, os.ModeNamedPipe)
		p := New(
			WithInput(f),
			WithFlagSet(flagSet, nil),
		)

		err := p.parseInput()
		require.NoError(t, err)
		require.Len(t, p.requestedFiles, 1)
		require.EqualValues(t, []string{"main.go"}, p.requestedFiles)
	})
}
