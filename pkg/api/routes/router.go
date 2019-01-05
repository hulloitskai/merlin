package routes

import (
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

// Router matches requests to the routes defined in this package.
type Router struct {
	*Config

	hr httprouter.Router
}

// NewRouter returns a new Router.
func NewRouter(cfg *Config) (*Router, error) {
	if cfg == nil {
		return nil, errors.New("routes: cannot create Router with a nil Config")
	}
	if cfg.Logger == nil {
		cfg.Logger = zap.NewNop().Sugar()
	}

	// Make and configure httprouter.Router.
	hr := httprouter.New()
	hr.RedirectTrailingSlash = false

	r := &Router{
		Config: cfg,
		hr:     *hr,
	}
	r.registerRoutes()
	return r, nil
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.hr.ServeHTTP(w, req)
}

func (r *Router) registerRoutes() {
	router := &r.hr
	registerIndex(router, r.Config.Logger.Named("index"))
	registerSheets(router, r.Config.Scraper, r.Config.Logger.Named("sheets"))
}
