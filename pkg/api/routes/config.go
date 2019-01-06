package routes

import (
	"go.uber.org/zap"

	"github.com/stevenxie/merlin/pkg/models"
	"github.com/stevenxie/merlin/pkg/models/balance"
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
	ScrapeFilings(ticker string) (*models.FilingResults, error)
}
