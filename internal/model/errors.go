package model

type errNoCoverageFile struct{ error }

func (e errNoCoverageFile) Title() string { return "Coverage report not found" }
func (e errNoCoverageFile) Description() string {
	return `Requested coverage profile is not found.
By default, "coverage.out" is used. For other names, use "--profile" flag.`
}
func (e errNoCoverageFile) OriginalError() error { return e }

type errInvalidCoverageFile struct{ error }

func (e errInvalidCoverageFile) Title() string { return "Invalid coverage file" }
func (e errInvalidCoverageFile) Description() string {
	return `The provided coverage file was found, but can't be parsed.
Update the coverage report and try again.`
}
func (e errInvalidCoverageFile) OriginalError() error { return e }

type errNoProfiles struct{ error }

func (e errNoProfiles) Title() string { return "No coverage data" }
func (e errNoProfiles) Description() string {
	return `Provided coverage profile is valid, but it doesn't have any entries.
Try adding tests and generating the report again.`
}
func (e errNoProfiles) OriginalError() error { return nil }

type errGoModNotFound struct{ error }

func (e errGoModNotFound) Title() string { return "go.mod file is not available" }
func (e errGoModNotFound) Description() string {
	return `This program must be executed from Go project root.
Run the program from the folder that contains go.mod file.`
}
func (e errGoModNotFound) OriginalError() error { return e }

type errInvalidGoMod struct{ error }

func (e errInvalidGoMod) Title() string { return "Invalid go.mod file" }
func (e errInvalidGoMod) Description() string {
	return "go.mod file does not include a valid Go module name: `module ...`"
}
func (e errInvalidGoMod) OriginalError() error { return nil }

type errSourceFileNotFound struct{ error }

func (e errSourceFileNotFound) Title() string { return "Source code file not found" }
func (e errSourceFileNotFound) Description() string {
	return `A file that appears in the coverage report was not found in the file tree.
Update the coverage report and try again.`
}
func (e errSourceFileNotFound) OriginalError() error { return e }

type errCantOpenSourceFile struct{ error }

func (e errCantOpenSourceFile) Title() string { return "Can't open source code file" }
func (e errCantOpenSourceFile) Description() string {
	return "Source code file was found, but cannot be opened."
}
func (e errCantOpenSourceFile) OriginalError() error { return e }

type errMismatchingProfile struct{ error }

func (e errMismatchingProfile) Title() string { return "Coverage data doesn't match the source" }
func (e errMismatchingProfile) Description() string {
	return `Coverage data cannot be applied to the existing source code.
Update the coverage report and try again.`
}
func (e errMismatchingProfile) OriginalError() error { return e }
