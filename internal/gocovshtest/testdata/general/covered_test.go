package general_test

import (
	"testing"

	"github.com/orlangure/gocovsh/internal/model/testdata/general"
	"github.com/stretchr/testify/require"
)

func TestFull(t *testing.T) {
	require.Equal(t, "full", general.Full())
}
