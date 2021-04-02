package page

var pages []*Page

func InsertPage(p *Page) {
	// If the page has no publish date, do not add it to the pages list.
	if p.Options.PublishDate == "" {
		return
	}

	index := 0
	for k, v := range pages {
		if p.Options.PublishDate > v.Options.PublishDate {
			index = k
			break
		}
		index++
	}

	pages = append(pages, nil)
	copy(pages[index+1:], pages[index:])
	pages[index] = p

}
