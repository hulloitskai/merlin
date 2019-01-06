package routes

import (
	"net/http"

	"go.uber.org/zap"

	hr "github.com/julienschmidt/httprouter"
	ess "github.com/unixpickle/essentials"
)

func registerNotes(r *hr.Router, s Scraper, logger *zap.SugaredLogger) {
	nh := &notesHandler{Scraper: s, l: logger}
	nh.RegisterTo(r)
}

type notesHandler struct {
	Scraper
	l *zap.SugaredLogger
}

func (nh *notesHandler) RegisterTo(r *hr.Router) {
	r.GET("/notes/:cik/:accNum", nh.Handle)
	r.GET("/notes/:cik/:accNum/", handleTrailingSlashRedir)
}

func (nh *notesHandler) Handle(w http.ResponseWriter, _ *http.Request,
	params hr.Params) {
	var (
		cik        = params.ByName("cik")
		accNum     = params.ByName("accNum")
		notes, err = nh.ScrapeFinanceNotes(cik, accNum)
		rw         = responseWriter{w}
	)

	if err != nil {
		nh.l.Debugf("Error while scraping finance notes for cik='%s', "+
			"accNum='%s': %v", cik, accNum, err)

		w.WriteHeader(http.StatusInternalServerError)
		ess.AddCtxTo("routes: scraping finance notes", &err)
		jerr := jsonErrorFrom(err)
		if err = rw.WriteJSON(&jerr); err != nil {
			nh.l.Errorf("Error writing JSON response: %v", err)
		}
		return
	}

	rw.WriteJSON(notes)
}
