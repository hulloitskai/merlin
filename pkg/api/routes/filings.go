package routes

import (
	"net/http"
	"strings"

	"go.uber.org/zap"

	hr "github.com/julienschmidt/httprouter"
	"github.com/stevenxie/merlin/pkg/models"
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
	r.GET("/filings/:ticker/latest/10k", fh.HandleLatest10K)
	r.GET("/filings/:ticker/latest/10k/", handleTrailingSlashRedir)
}

func (fh *filingsHandler) Handle(w http.ResponseWriter, r *http.Request,
	params hr.Params) {
	var (
		rw      = responseWriter{w, fh.l}
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
		rw      = responseWriter{w, fh.l}
		res, ok = fh.partiallyHandle(&rw, r, params)
	)
	if !ok {
		return
	}

	// Find first filing of type 10-K.
	var f10k struct {
		CIK           string
		models.Filing `json:"filing"`
	}
	f10k.CIK = res.CIK
	for _, filing := range res.Filings {
		if strings.ToLower(filing.Type) == "10-k" {
			f10k.Filing = *filing
		}
	}
	rw.WriteJSON(&f10k)
}

func (fh *filingsHandler) partiallyHandle(rw *responseWriter, _ *http.Request,
	params hr.Params) (res *models.FilingResults, ok bool) {
	ticker := params.ByName("ticker")
	res, err := fh.ScrapeFilings(ticker)

	if err != nil {
		fh.l.Debugf("Error while scraping company filings for ticker='%s': %v",
			ticker, err)
		ess.AddCtxTo("routes: scraping company filings", &err)

		// Check if is a Content-Type error; if it is, then we assume that an
		// invalid ticker was provided.
		var (
			code = http.StatusInternalServerError
			desc string
		)
		if strings.Contains(err.Error(), "Content-Type") {
			code = http.StatusBadRequest
			desc = "Provided ticker is likely invalid."
		}

		rw.WriteHeader(code)
		jerr := jsonErrorFrom(err, code)
		if desc != "" {
			jerr.Desc = desc
		}

		rw.WriteJSON(&jerr)
		return nil, false
	}

	return res, true
}
