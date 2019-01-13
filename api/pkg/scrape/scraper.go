package scrape

import (
	"net/http"

	gq "github.com/PuerkitoBio/goquery"
	ess "github.com/unixpickle/essentials"
)

// A Scraper is capable of scraping filings from EDGAR.
type Scraper struct {
	Client *http.Client
}

// NewScraper returns a new Scraper
func NewScraper() *Scraper {
	return &Scraper{Client: new(http.Client)}
}

func (s *Scraper) readDocumentAt(url string) (*gq.Document, error) {
	res, err := s.Client.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	doc, err := gq.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, ess.AddCtx("parsing response body with goquery", err)
	}

	err = res.Body.Close()
	return doc, ess.AddCtx("closing response body", err)
}
