package internal

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {

	t.Skip("infra")
	// env.Initialize()

	template, err := Run()
	require.NoError(t, err, template)
}
