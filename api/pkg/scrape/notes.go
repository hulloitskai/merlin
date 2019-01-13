package scrape

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	gq "github.com/PuerkitoBio/goquery"
	"github.com/stevenxie/merlin/api/pkg/models"
	ess "github.com/unixpickle/essentials"
)

// ScrapeFinanceNotes scrapes EDGAR for the balance sheets for the filing
// associated with cik and accNum.
func (s *Scraper) ScrapeFinanceNotes(cik, accNum string) (
	models.Notes, error) {
	var (
		url      = makeViewerURL(cik, accNum)
		doc, err = s.readDocumentAt(url)
	)
	if err != nil {
		return nil, err
	}

	sel := doc.Find("#menu_cat3")
	if sel.Length() != 1 {
		return nil, errors.New("scrape: could not find notes sidebar item")
	}
	sel = sel.Parent().Find("ul")
	if sel.Length() != 1 {
		return nil, errors.New("scrape: could not find list of note items")
	}

	var links []string
	sel.Children().EachWithBreak(func(i int, row *gq.Selection) bool {
		id, ok := row.Attr("id")
		if !ok {
			err = fmt.Errorf("row %d has no 'id' attribute", i)
			return false
		}

		id = strings.ToUpper(id)
		var (
			flatAccNum = strings.Replace(accNum, "-", "", -1)
			link       = fmt.Sprintf("%s/%s/%s/%s.htm", edgarDataURL, cik, flatAccNum,
				id)
		)
		links = append(links, link)
		return true
	})

	// Parse note at each link.
	notes := make(models.Notes, 0)
	for _, link := range links {
		note, err := s.scrapeFinanceNoteAt(link)
		if err != nil {
			return nil, ess.AddCtx(fmt.Sprintf("scraping note at '%s'", link), err)
		}

		if note == nil {
			continue
		}
		notes = append(notes, note)
	}
	return notes, nil
}

func (s *Scraper) scrapeFinanceNoteAt(url string) (*models.Note, error) {
	doc, err := s.readDocumentAt(url)
	if err != nil {
		return nil, err
	}

	var (
		sel    = doc.Find(".text")
		expect = 2
	)
	if length := sel.Length(); length < expect {
		return nil, fmt.Errorf("did not find expected number of '.text' elements "+
			"(expected %d, but got %d)", expect, length)
	}

	// Strip first '.text' element (empty section header artifact).
	sel.Nodes = sel.Nodes[1:]

	sel = sel.First().Children().First().Children().First()
	if sel.Length() != 1 {
		return nil, errors.New("could not find note title element")
	}

	// Trim note tag from title.
	title := strings.TrimSpace(sel.Text())
	const marker = "Note "
	if !strings.HasPrefix(title, marker) { // not a real note
		return nil, nil
	}

	title = strings.TrimPrefix(title, marker)
	closei := strings.IndexByte(title, ' ')
	if closei == -1 {
		return nil, errors.New("bad spacing in note tag (pre)")
	}

	sid := title[:closei]
	title = strings.TrimPrefix(title[closei+1:], "– ")
	id, err := strconv.Atoi(sid)
	if err != nil {
		return nil, ess.AddCtx("parsing note ID as integer", err)
	}

	return &models.Note{
		ID:    id,
		Title: title,
		Link:  url,
	}, nil
}
