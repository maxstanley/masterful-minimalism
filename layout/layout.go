package layout

import (
	"io/ioutil"
	"path"
	"strings"
)

type Layout struct {
	Name     string
	Template string
}

func NewFromFile(filePath string) (layout Layout, err error) {
	layout.Name = strings.ToLower(strings.TrimSuffix(path.Base(filePath), ".html"))
	var file []byte

	if file, err = ioutil.ReadFile(filePath); err != nil {
		return
	}

	layout.Template = string(file)
	return
}
