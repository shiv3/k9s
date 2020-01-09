package config_test

import (
	"testing"

	"github.com/derailed/k9s/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestPluginLoad(t *testing.T) {
	p := config.NewPlugins()
	assert.Nil(t, p.LoadPlugins("test_assets/plugin.yml"))

	assert.Equal(t, 1, len(p.Plugin))
	k, ok := p.Plugin["blah"]
	assert.True(t, ok)
	assert.Equal(t, "shift-s", k.ShortCut)
	assert.Equal(t, "blee", k.Description)
	assert.Equal(t, []string{"po", "dp"}, k.Scopes)
	assert.Equal(t, "duh", k.Command)
	assert.Equal(t, []string{"-n", "$NAMESPACE", "-boolean"}, k.Args)
}
