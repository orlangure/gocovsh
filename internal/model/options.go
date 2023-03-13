package model

// Option is a function that can be used to modify the model.
type Option func(*Model)

// WithProfileFilename sets the filename of the coverage report to be loaded.
func WithProfileFilename(name string) Option {
	return func(m *Model) {
		m.profileFilename = name
	}
}

// WithCoverageSorting asks for the profiles to be sorted by coverage percent instead of alphabetically.
func WithCoverageSorting(sortByCoverage bool) Option {
	return func(m *Model) {
        if sortByCoverage {
            m.sortState.Type = sortStateByPercentage
        }
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

// WithFilteredLines sets a list of lines to display for every file. Other
// lines will not appear.
func WithFilteredLines(files map[string][]int) Option {
	return func(m *Model) {
		m.filteredLinesByFile = make(map[string][]int, len(files))
		uniqueLines := map[int]interface{}{}

		for file, lines := range files {
			linesWithContext := make([]int, 0, len(lines))

			for _, line := range lines {
				if _, ok := uniqueLines[line]; !ok {
					linesWithContext = append(linesWithContext, line)
					uniqueLines[line] = nil
				}
			}

			m.filteredLinesByFile[file] = linesWithContext
		}
	}
}
