package routes

import (
	"net/http"

	"go.uber.org/zap"

	hr "github.com/julienschmidt/httprouter"
	ess "github.com/unixpickle/essentials"
)

func registerSheets(r *hr.Router, s Scraper, logger *zap.SugaredLogger) {
	sh := &sheetsHandler{Scraper: s, l: logger}
	sh.RegisterTo(r)
}

type sheetsHandler struct {
	Scraper
	l *zap.SugaredLogger
}

func (sh *sheetsHandler) RegisterTo(r *hr.Router) {
	r.GET("/sheets/:cik/:accNum", sh.Handle)
	r.GET("/sheets/:cik/:accNum/", handleTrailingSlashRedir)
}

func (sh *sheetsHandler) Handle(w http.ResponseWriter, _ *http.Request,
	params hr.Params) {
	var (
		cik         = params.ByName("cik")
		accNum      = params.ByName("accNum")
		sheets, err = sh.ScrapeBalanceSheets(cik, accNum)
		rw          = responseWriter{w, sh.l}
	)

	if err != nil {
		sh.l.Debugf("Error while scraping balance sheets for cik='%s', "+
			"accNum='%s': %v", cik, accNum, err)
		ess.AddCtxTo("routes: scraping balance sheets", &err)

		code := http.StatusInternalServerError
		w.WriteHeader(code)
		jerr := jsonErrorFrom(err, code)
		rw.WriteJSON(&jerr)
	}

	rw.WriteJSON(sheets)
}
