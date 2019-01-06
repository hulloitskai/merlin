package routes

import (
	"net/http"
	"strings"

	"github.com/stevenxie/merlin/pkg/models"

	"go.uber.org/zap"

	hr "github.com/julienschmidt/httprouter"
	ess "github.com/unixpickle/essentials"
)

func registerFilings(r *hr.Router, s Scraper, logger *zap.SugaredLogger) {
	fh := &filingsHandler{Scraper: s, l: logger}
	fh.RegisterTo(r)
}

type filingsHandler struct {
	Scraper
	l *zap.SugaredLogger
}

func (fh *filingsHandler) RegisterTo(r *hr.Router) {
	r.GET("/filings/:ticker", handleTrailingSlashRedir)
	r.GET("/filings/:ticker/", fh.Handle)
	r.GET("/filings/:ticker/latest10k", fh.HandleLatest10K)
	r.GET("/filings/:ticker/latest10k/", handleTrailingSlashRedir)
}

func (fh *filingsHandler) Handle(w http.ResponseWriter, r *http.Request,
	params hr.Params) {
	var (
		rw      = responseWriter{w}
		res, ok = fh.partiallyHandle(&rw, r, params)
	)
	if !ok {
		return
	}
	rw.WriteJSON(res)
}

func (fh *filingsHandler) HandleLatest10K(w http.ResponseWriter,
	r *http.Request, params hr.Params) {
	var (
		rw      = responseWriter{w}
		res, ok = fh.partiallyHandle(&rw, r, params)
	)
	if !ok {
		return
	}

	// Find first filing of type 10-K.
	var f10k *models.Filing
	for _, filing := range res.Filings {
		if strings.ToLower(filing.Type) == "10-k" {
			f10k = filing
		}
	}
	rw.WriteJSON(f10k)
}

func (fh *filingsHandler) partiallyHandle(rw *responseWriter, _ *http.Request,
	params hr.Params) (res *models.FilingResults, ok bool) {
	ticker := params.ByName("ticker")
	res, err := fh.ScrapeFilings(ticker)

	if err != nil {
		fh.l.Debugf("Error while scraping company filings for ticker='%s': %v",
			ticker, err)

		rw.WriteHeader(http.StatusInternalServerError)
		ess.AddCtxTo("routes: scraping company filings", &err)
		jerr := jsonErrorFrom(err)
		if err = rw.WriteJSON(&jerr); err != nil {
			fh.l.Errorf("Error writing JSON response: %v", err)
		}
		return nil, false
	}

	return res, true
}
