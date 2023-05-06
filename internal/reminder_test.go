package internal

import (
	"testing"

	"github.com/democracy-tools/countmein-scheduler/internal/ds"
	"github.com/democracy-tools/go-common/env"
	"github.com/democracy-tools/go-common/whatsapp"
	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {

	t.Skip("infra")
	// env.Initialize()

	template, err := Run(ds.NewClientWrapper(env.Project), whatsapp.NewClientWrapper())
	require.NoError(t, err, template)
}
