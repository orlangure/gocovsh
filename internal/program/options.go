package program

import (
	"flag"
	"io"
	"io/fs"
	"runtime/debug"
)

// Option is a function that can be passed to WithOptions.
type Option func(*Program)

// WithBuildInfo sets the version information of this program.
func WithBuildInfo(version, commit, date string) Option {
	return func(p *Program) {
		p.version = version
		p.commit = commit
		p.date = date
	}
}

// WithGoModInfo reads build info and sets its fields, if its available.
func WithGoModInfo() Option {
	return func(p *Program) {
		if info, ok := debug.ReadBuildInfo(); ok && info.Main.Sum != "" {
			p.modVersion = info.Main.Version
			p.modSum = info.Main.Sum
		}
	}
}

// WithLogFile sets the path to the log file.
func WithLogFile(path string) Option {
	return func(p *Program) {
		p.logFile = path
	}
}

// WithFlagSet is an optional way to set the flag set for the program. Useful
// in tests.
func WithFlagSet(fs *flag.FlagSet, args []string) Option {
	return func(p *Program) {
		p.flagSet = fs
		p.args = args

		if p.args == nil {
			p.args = []string{}
		}
	}
}

// WithOutput sets the stdout writer for the program. This should be used
// for testing some of the basic outputs (non-interactive).
func WithOutput(w io.Writer) Option {
	return func(p *Program) {
		p.output = w
	}
}

// WithInput sets the stdin for the program. This should be used for testing
// the features that read from stdin.
func WithInput(file fs.File) Option {
	return func(p *Program) {
		p.input = file
	}
}
