package errors_test

import (
	"testing"

	"github.com/orlangure/gocovsh/internal/model/testdata/errors"
	"github.com/stretchr/testify/require"
)

func TestFull(t *testing.T) {
	require.Equal(t, "full", errors.Full())
}
