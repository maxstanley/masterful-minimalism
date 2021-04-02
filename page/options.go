package page

import "github.com/maxstanley/masterful-minimalism/configuration"

type Options struct {
	Abstract    string   `yaml:"abstract"`
	AuthorList  []string `yaml:"authorList"`
	Layout      string   `yaml:"layout"`
	PublishDate string   `yaml:"publishDate"`
	LastUpdated string   `yaml:"lastUpdated"`
	Tags        []string `yaml:"tags"`
	Title       string   `yaml:"title"`

	Authors []configuration.Author
}
