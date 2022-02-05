package general_test

import (
	"testing"

	"github.com/orlangure/gocovsh/internal/model/testdata/general"
	"github.com/stretchr/testify/require"
)

func TestCovered(t *testing.T) {
	require.Equal(t, "covered", general.Covered())
}

func TestSecondCovered(t *testing.T) {
	require.Equal(t, "covered", general.SecondCovered())
}
