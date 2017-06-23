package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefaultConfig(t *testing.T) {
	c := Default()
	require.Equal(t, ".", c.Source)
	require.Equal(t, "./_site", c.Destination)
	require.Equal(t, "_layouts", c.LayoutsDir)
}

func TestUnmarshal(t *testing.T) {
	c := Default()
	Unmarshal([]byte(`source: x`), &c)
	require.Equal(t, "x", c.Source)
	require.Equal(t, "./_site", c.Destination)
}