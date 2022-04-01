package model

// Option is a function that can be used to modify the model.
type Option func(*Model)

// WithProfileFilename sets the filename of the coverage report to be loaded.
func WithProfileFilename(name string) Option {
	return func(m *Model) {
		m.profileFilename = name
	}
}

// WithCodeRoot sets the root directory of the code to be analyzed.
func WithCodeRoot(root string) Option {
	return func(m *Model) {
		m.codeRoot = root
	}
}

// WithRequestedFiles sets the list of files to be displayed.
func WithRequestedFiles(files []string) Option {
	return func(m *Model) {
		m.requestedFiles = make(map[string]bool)

		for _, v := range files {
			m.requestedFiles[v] = true
		}
	}
}
