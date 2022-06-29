package codeview

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestContextifyFilteredLines(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input           []int
		expectedActual  []int
		expectedContext map[int]bool
	}{
		{
			input:           []int{},
			expectedActual:  []int{},
			expectedContext: map[int]bool{},
		},
		{
			input:           []int{1, 2},
			expectedActual:  []int{1, 2, 3},
			expectedContext: map[int]bool{3: true},
		},
		{
			input:           []int{2, 3},
			expectedActual:  []int{1, 2, 3, 4},
			expectedContext: map[int]bool{1: true, 4: true},
		},
		{
			input:          []int{2, 3, 20},
			expectedActual: []int{1, 2, 3, 4, 19, 20, 21},
			expectedContext: map[int]bool{
				1: true, 4: true,
				19: true, 21: true,
			},
		},
		{
			input:          []int{2, 3, 20, 21},
			expectedActual: []int{1, 2, 3, 4, 19, 20, 21, 22},
			expectedContext: map[int]bool{
				1: true, 4: true,
				19: true, 22: true,
			},
		},
		{
			input:          []int{2, 3, 20, 30},
			expectedActual: []int{1, 2, 3, 4, 19, 20, 21, 29, 30, 31},
			expectedContext: map[int]bool{
				1: true, 4: true,
				19: true, 21: true,
				29: true, 31: true,
			},
		},
		{
			input:          []int{2, 3, 20, 22},
			expectedActual: []int{1, 2, 3, 4, 19, 20, 21, 22, 23},
			expectedContext: map[int]bool{
				1: true, 4: true,
				19: true, 21: true, 23: true,
			},
		},
		{
			input:          []int{2, 3, 20, 21, 22, 24, 25, 30},
			expectedActual: []int{1, 2, 3, 4, 19, 20, 21, 22, 23, 24, 25, 26, 29, 30, 31},
			expectedContext: map[int]bool{
				1: true, 4: true,
				19: true, 23: true,
				26: true,
				29: true, 31: true,
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%v", test.input), func(t *testing.T) {
			output := contextifyFilteredLines(test.input)
			require.EqualValues(t, test.expectedActual, output.actualLines)
			require.EqualValues(t, test.expectedContext, output.contextLines)
		})
	}
}

func Range(from, to int) []int {
	result := make([]int, 0, to-from)

	for i := 0; i <= to; i++ {
		result = append(result, i)
	}

	return result
}
