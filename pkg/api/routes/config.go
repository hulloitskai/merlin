package routes

import (
	"github.com/stevenxie/merlin/pkg/models"
	"github.com/stevenxie/merlin/pkg/models/balance"
	"go.uber.org/zap"
)

// A Config is used to configure a Router.
type Config struct {
	Scraper Scraper
	Logger  *zap.SugaredLogger
}

// A Scraper is capable of scraping filings from EDGAR.
type Scraper interface {
	ScrapeBalanceSheets(cik, accNum string) (balance.Sheets, error)
	ScrapeFinanceNotes(cik, accNum string) (models.Notes, error)
}
