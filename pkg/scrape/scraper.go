package scrape

import (
	"net/http"
)

// A Scraper is capable of scraping filings from EDGAR.
type Scraper struct {
	Client *http.Client
}

// NewScraper returns a new Scraper
func NewScraper() *Scraper {
	return &Scraper{Client: new(http.Client)}
}
