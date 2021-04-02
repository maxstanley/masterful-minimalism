package configuration

import (
	"io/ioutil"

	"github.com/maxstanley/masterful-minimalism/layout"
	"github.com/snabb/sitemap"
	"gopkg.in/yaml.v2"
)

type Configuration struct {
	// Authors contains information about the sites authors.
	Authors map[string]Author `yaml:"authors"`
	// Site contains the configuration information for the site.
	Site Site `yaml:"site"`

	// Command Line Options
	LayoutsPath string
	PagesPath   string
	OutputPath  string

	Layouts map[string]layout.Layout

	SiteMap *sitemap.SitemapIndex

	// Other contains custom configuration options.
	Other map[string]interface{} `yaml:"-"`
}

func NewFromFile(configurationPath string) (configuration *Configuration, err error) {
	configuration = &Configuration{}
	configuration.Other = map[string]interface{}{}

	var file []byte

	if file, err = ioutil.ReadFile(configurationPath); err != nil {
		return
	}

	if err = yaml.Unmarshal(file, configuration); err != nil {
		return
	}

	if err = yaml.Unmarshal(file, configuration.Other); err != nil {
		return
	}

	configuration.SiteMap = sitemap.NewSitemapIndex()

	return
}

func (c *Configuration) AddPaths(layouts string, pages string, output string) {
	c.LayoutsPath = layouts
	c.PagesPath = pages
	c.OutputPath = output
}
