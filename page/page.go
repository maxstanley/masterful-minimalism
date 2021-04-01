package page

import (
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	handlebars "github.com/aymerick/raymond"
	"github.com/russross/blackfriday/v2"
	"gopkg.in/yaml.v2"

	"github.com/maxstanley/masterful-minimalism/configuration"
	"github.com/maxstanley/masterful-minimalism/layout"
)

type Heading struct {
	Value    int
	Text     string
	IDLink   string
	Parent   *Heading
	Children []*Heading
}

type Page struct {
	outputPath string

	Options  Options                      `handlebars:"document"`
	Content  string                       `handlebars:"content"`
	Headings []*Heading                   `handlebars:"headings"`
	Global   *configuration.Configuration `handlebars:"global"`
	Pages    []*Page                      `handlebars:"pages"`

	RelativePath string
}

func NewFromFile(pagePath string, configurationOptions *configuration.Configuration) (page *Page, err error) {
	page = &Page{Global: configurationOptions, Pages: pages}

	var (
		file []byte
	)

	if file, err = ioutil.ReadFile(pagePath); err != nil {
		return
	}

	fileOutputPath := strings.Replace(pagePath, configurationOptions.PagesPath, configurationOptions.OutputPath, 1)
	fileOutputPath = strings.Replace(fileOutputPath, ".md", ".html", 1)
	page.outputPath = fileOutputPath

	page.RelativePath = strings.Replace(fileOutputPath, configurationOptions.OutputPath, "", 1)

	if err = yaml.Unmarshal(file, &page.Options); err != nil {
		return
	}

	md := []byte(strings.Trim(strings.SplitN(string(file), "---", 3)[2], "\n"))

	extensions := blackfriday.Tables | blackfriday.FencedCode | blackfriday.AutoHeadingIDs | blackfriday.Strikethrough | blackfriday.SpaceHeadings | blackfriday.Footnotes | blackfriday.DefinitionLists

	page.Content = string(blackfriday.Run(
		md,
		blackfriday.WithNoExtensions(),
		blackfriday.WithExtensions(extensions),
	))

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(page.Content))
	if err != nil {
		return
	}

	headingSelector := doc.Find("h1,h2,h3,h4,h5,h6") //.Find("h2").Find("h3").Find("h4").Find("h5").Find("h6")
	headings := []*Heading{}
	var previousHeading *Heading = nil
	headingSelector.Each(func(_ int, sel *goquery.Selection) {
		value, err := strconv.Atoi(sel.Nodes[0].Data[1:2])
		if err != nil {
			return
		}

		h := &Heading{
			Value:  value,
			Text:   sel.Text(),
			IDLink: blackfriday.SanitizedAnchorName(sel.Text()),
			Parent: nil,
		}

		if value == 1 {
			// if h1 add to headings array
			headings = append(headings, h)
		} else if value > previousHeading.Value {
			// if the heading is a sub heading of the previous
			h.Parent = previousHeading
			previousHeading.Children = append(previousHeading.Children, h)
		} else {
			// if the heading is a sub heading of a parent of the previous heading
			for value <= previousHeading.Value {
				previousHeading = previousHeading.Parent
			}
			h.Parent = previousHeading
			previousHeading.Children = append(previousHeading.Children, h)
		}

		previousHeading = h
	})
	page.Headings = headings

	for _, authorString := range page.Options.AuthorList {
		if author, ok := configuration.Authors[authorString]; ok {
			page.Options.Authors = append(page.Options.Authors, author)
		}
	}

	InsertPage(page)
	return
}

func (p *Page) CreateFile(configurationOptions *configuration.Configuration, l layout.Layout) error {

	outputFile := handlebars.MustRender(l.Template, p)
	outputFile = handlebars.MustRender(outputFile, p)

	directory, _ := path.Split(p.outputPath)
	os.MkdirAll(directory, 0700)

	if err := ioutil.WriteFile(p.outputPath, []byte(outputFile), 0644); err != nil {
		return err
	}

	return nil
}
