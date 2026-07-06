package daemon

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunArguments(t *testing.T) {
	assert.Equal(t, []string{RunCommand}, runArguments(""))
	assert.Equal(t, []string{RunCommand, "--env-file", "/etc/pihole-sync/.env"}, runArguments("/etc/pihole-sync/.env"))
}

func TestNew(t *testing.T) {
	s, err := New("")
	require.NoError(t, err)
	require.NotNil(t, s)
}
