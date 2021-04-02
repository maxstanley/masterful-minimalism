package main

import (
	"flag"
	"log"
	"os"
	"path"

	"github.com/maxstanley/masterful-minimalism/configuration"
	"github.com/maxstanley/masterful-minimalism/layout"
	"github.com/maxstanley/masterful-minimalism/page"
	"github.com/maxstanley/masterful-minimalism/utils"
)

func main() {
	var (
		workingDirectory  *string = new(string)
		configurationPath *string = new(string)
		layoutsPath       *string = new(string)
		pagesPath         *string = new(string)
		outputPath        *string = new(string)
	)

	flag.StringVar(workingDirectory, "input", ".", "Path to the working directory for the execution.")
	flag.StringVar(configurationPath, "config", "config.yaml", "Path to configuration file.")
	flag.StringVar(layoutsPath, "layouts", "layouts/", "Path to layouts folder.")
	flag.StringVar(pagesPath, "pages", "pages/", "Path to pages folder.")
	flag.StringVar(outputPath, "output", "output/", "Path to output folder.")
	flag.Parse()

	dirVariables := []*string{workingDirectory, outputPath}
	for _, dir := range dirVariables {
		// If the input working directory is not absolute.
		if !path.IsAbs(*dir) {
			// Get the current working directory.
			cwd, err := os.Getwd()
			if err != nil {
				log.Fatalln(err.Error())
			}
			// Create an absolute path from the input.
			*dir = path.Join(cwd, *dir)
		}
	}

	directoryVariables := []*string{configurationPath, layoutsPath, pagesPath}
	for _, dir := range directoryVariables {
		if !path.IsAbs(*dir) {
			*dir = path.Join(*workingDirectory, *dir)
		}
	}

	os.RemoveAll(*outputPath)
	os.MkdirAll(*outputPath, 0700)

	var (
		configurationOptions *configuration.Configuration
		err                  error
	)

	if configurationOptions, err = configuration.NewFromFile(*configurationPath); err != nil {
		log.Fatalln(err.Error())
	}

	configurationOptions.AddPaths(*layoutsPath, *pagesPath, *outputPath)

	configuration.Authors = configurationOptions.Authors

	layoutPaths, err := utils.WalkFolderFiles(*layoutsPath)
	if err != nil {
		log.Fatalln(err.Error())
	}

	layouts := map[string]layout.Layout{}
	for _, layoutPath := range layoutPaths {
		var l layout.Layout
		if l, err = layout.NewFromFile(layoutPath); err != nil {
			log.Fatalln(err.Error())
		}
		layouts[l.Name] = l
	}

	configurationOptions.Layouts = layouts

	pagePaths, err := utils.WalkFolderFiles(*pagesPath)
	if err != nil {
		log.Fatalln(err.Error())
	}

	var pages []*page.Page
	for _, pagePath := range pagePaths {
		p, err := page.NewFromFile(pagePath, configurationOptions)
		if err != nil {
			log.Fatalln(err.Error())
		}
		pages = append(pages, p)
	}

	for _, p := range pages {
		l := layouts[p.Options.Layout]
		p.CreateFile(configurationOptions, l)
	}
}
