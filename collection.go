package gojekyll

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Collection is a Jekyll collection.
type Collection struct {
	Site   *Site
	Name   string
	Data   VariableMap
	Output bool
	Pages  []Page
}

// NewCollection creates a new Collection with defaults d
func NewCollection(s *Site, name string, d VariableMap) *Collection {
	return &Collection{
		Site:   s,
		Name:   name,
		Data:   d,
		Output: d.Bool("output", false),
	}
}

// ReadCollections reads the pages of the collections named in the site configuration.
// It adds each collection's pages to the site map, and creates a template site variable for each collection.
func (s *Site) ReadCollections() error {
	for name, d := range s.config.Collections {
		c := NewCollection(s, name, d)
		s.Collections = append(s.Collections, c)
		if c.Output { // TODO always read the pages; just don't build them / include them in routes
			if err := c.ReadPages(); err != nil {
				return err
			}
		}
	}
	return nil
}

// PageTemplateObjects returns an array of page objects, for use as the template variable
// value of the collection.
func (c *Collection) PageTemplateObjects() (d []VariableMap) {
	for _, p := range c.Pages {
		d = append(d, p.TemplateObject())
	}
	return
}

// IsPosts returns true if the collection is the special "posts" collection.
func (c *Collection) IsPosts() bool {
	return c.Name == "posts"
}

// PathPrefix returns the collection's directory prefix, e.g. "_posts/"
func (c *Collection) PathPrefix() string {
	return filepath.FromSlash("_" + c.Name + "/")
}

// Source returns the source directory for pages in the collection.
func (c *Collection) Source() string {
	return filepath.Join(c.Site.Source, "_"+c.Name)
}

// ReadPages scans the file system for collection pages, and adds them to c.Pages.
func (c *Collection) ReadPages() error {
	basePath := c.Site.Source
	collectionDefaults := MergeVariableMaps(c.Data, VariableMap{
		"collection": c.Name,
	})

	walkFn := func(name string, info os.FileInfo, err error) error {
		if err != nil {
			// if the issue is simply that the directory doesn't exist, warn instead of error
			if os.IsNotExist(err) {
				if !c.IsPosts() {
					fmt.Printf("Missing directory for collection: _%s\n", c.Name)
				}
				return nil
			}
			return err
		}
		relname, err := filepath.Rel(basePath, name)
		switch {
		case strings.HasPrefix(filepath.Base(relname), "."):
			return nil
		case err != nil:
			return err
		case info.IsDir():
			return nil
		}
		defaults := MergeVariableMaps(c.Site.GetFrontMatterDefaults(relname, ""), collectionDefaults)
		p, err := ReadPage(c.Site, c, relname, defaults)
		switch {
		case err != nil:
			return err
		case p.Static():
			fmt.Printf("skipping static file inside collection: %s\n", name)
		case p.Published():
			c.Site.Paths[p.Permalink()] = p
			c.Pages = append(c.Pages, p)
		}
		return nil
	}
	return filepath.Walk(c.Source(), walkFn)
}
