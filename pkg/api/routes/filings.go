package routes

import (
	"net/http"

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
	r.GET("/filings/:ticker", fh.Handle)
	r.GET("/filings/:ticker/", handleTrailingSlashRedir)
}

func (fh *filingsHandler) Handle(w http.ResponseWriter, _ *http.Request,
	params hr.Params) {
	var (
		ticker     = params.ByName("ticker")
		notes, err = fh.ScrapeFilings(ticker)
		rw         = responseWriter{w}
	)

	if err != nil {
		fh.l.Debugf("Error while scraping company filings for ticker='%s': %v",
			ticker, err)

		w.WriteHeader(http.StatusInternalServerError)
		ess.AddCtxTo("routes: scraping company filings", &err)
		jerr := jsonErrorFrom(err)
		if err = rw.WriteJSON(&jerr); err != nil {
			fh.l.Errorf("Error writing JSON response: %v", err)
		}
		return
	}

	rw.WriteJSON(notes)
}
