package plugins

import (
	"regexp"
	"testing"

	"github.com/osteele/gojekyll/config"
	"github.com/osteele/gojekyll/pages"
	"github.com/osteele/liquid"
	"github.com/stretchr/testify/require"
)

type siteFake struct {
	c config.Config
	e liquid.Engine
}

func (s siteFake) AddDocument(pages.Document, bool) {}
func (s siteFake) Config() *config.Config           { return &s.c }
func (s siteFake) Pages() (ps []pages.Page)         { return }

func TestAvatarTag(t *testing.T) {
	engine := liquid.NewEngine()
	plugins := []string{"jekyll-avatar"}
	Install(plugins, siteFake{config.Default(), engine})
	require.NoError(t, directory[plugins[0]].ConfigureTemplateEngine(engine))
	bindings := map[string]interface{}{"user": "osteele"}

	s, err := engine.ParseAndRenderString(`{% avatar osteele %}`, bindings)
	require.NoError(t, err)
	re := regexp.MustCompile(`<img class="avatar.*avatar.*usercontent\.com/osteele\b`)
	require.True(t, re.MatchString(s))

	s, err = engine.ParseAndRenderString(`{% avatar user='osteele' %}`, bindings)
	require.NoError(t, err)
	require.True(t, re.MatchString(s))

	s, err = engine.ParseAndRenderString(`{% avatar user=user %}`, bindings)
	require.NoError(t, err)
	require.True(t, re.MatchString(s))

	s, err = engine.ParseAndRenderString(`{% avatar user=user size=20 %}`, bindings)
	require.NoError(t, err)
	require.Contains(t, s, "20")
}
