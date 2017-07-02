package plugins

import (
	"fmt"
	"io"

	"github.com/osteele/gojekyll/tags"
	"github.com/osteele/liquid/chunks"
)

func init() {
	registerPlugin("jekyll-gist", func(_ PluginContext, h pluginHelper) error {
		h.tag("gist", gistTag)
		return nil
	})
}

func gistTag(w io.Writer, ctx chunks.RenderContext) error {
	argsline, err := ctx.ParseTagArgs()
	if err != nil {
		return err
	}
	args, err := tags.ParseArgs(argsline)
	if err != nil {
		return err
	}
	if len(args.Args) < 1 {
		return fmt.Errorf("gist tag: missing argument")
	}
	url := fmt.Sprintf("https://gist.github.com/%s.js", args.Args[0])
	if len(args.Args) >= 2 {
		url += fmt.Sprintf("?file=%s", args.Args[1])
	}
	_, err = w.Write([]byte(`<script src=` + url + `> </script>`))
	return err
}