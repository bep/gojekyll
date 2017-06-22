package pages

import (
	"testing"

	"github.com/osteele/gojekyll/templates"
	"github.com/stretchr/testify/require"
)

type containerMock struct{ pathPrefix string }

func (c containerMock) Output() bool       { return true }
func (c containerMock) PathPrefix() string { return c.pathPrefix }

func TestExpandPermalinkPattern(t *testing.T) {
	var (
		c    = containerMock{}
		d    = templates.VariableMap{}
		path = "/a/b/base.html"
	)

	testPermalinkPattern := func(pattern, path string, data templates.VariableMap) (string, error) {
		vs := templates.MergeVariableMaps(data, templates.VariableMap{"permalink": pattern})
		p := pageFields{container: c, relpath: path, frontMatter: vs}
		return p.expandPermalink()
	}

	// t.Run(":output_ext", func(t *testing.T) {
	// 	p, err := testPermalinkPattern("/base:output_ext", path, d)
	// 	require.NoError(t, err)
	// 	require.Equal(t, "/base.html", p)
	// })
	// t.Run(":output_ext renames markdown to .html", func(t *testing.T) {
	// 	p, err := testPermalinkPattern("/base:output_ext", "/a/b/base.md", d)
	// 	require.NoError(t, err)
	// 	require.Equal(t, "/base.html", p)
	// 	p, err = testPermalinkPattern("/base:output_ext", "/a/b/base.markdown", d)
	// 	require.NoError(t, err)
	// 	require.Equal(t, "/base.html", p)
	// })
	t.Run(":name", func(t *testing.T) {
		p, err := testPermalinkPattern("/name/:name", path, d)
		require.NoError(t, err)
		require.Equal(t, "/name/base", p)
	})
	t.Run(":path", func(t *testing.T) {
		p, err := testPermalinkPattern("/prefix:path/post", path, d)
		require.NoError(t, err)
		require.Equal(t, "/prefix/a/b/base/post", p)
	})
	t.Run(":title", func(t *testing.T) {
		p, err := testPermalinkPattern("/title/:title.html", path, d)
		require.NoError(t, err)
		require.Equal(t, "/title/base.html", p)
	})
	t.Run("invalid template variable", func(t *testing.T) {
		_, err := testPermalinkPattern("/:invalid", path, d)
		require.Error(t, err)
	})

	c = containerMock{"_c/"}
	path = "_c/a/b/c.d"
	t.Run(":path", func(t *testing.T) {
		p, err := testPermalinkPattern("/prefix:path/post", path, d)
		require.NoError(t, err)
		require.Equal(t, "/prefix/a/b/c/post", p)
	})
}