package routes

import (
	"fmt"
	"net/http"
	"os"

	"go.uber.org/zap"

	hr "github.com/julienschmidt/httprouter"
	"github.com/stevenxie/merlin/api/internal/info"
)

func registerIndex(r *hr.Router, logger *zap.SugaredLogger) {
	ih := newIndexHandler(logger)
	ih.RegisterTo(r)
}

type apiInfo struct {
	ID      string `json:"id"`
	Version string `json:"version"`
	Env     string `json:"environment"`
}

type indexHandler struct {
	template apiInfo
	l        *zap.SugaredLogger
}

func newIndexHandler(logger *zap.SugaredLogger) *indexHandler {
	ai := apiInfo{Version: info.Version}
	if ai.Version != "unset" {
		ai.ID = fmt.Sprintf("%s-api-%s", info.Namespace, info.Version)
	} else {
		ai.ID = fmt.Sprintf("%s-api", info.Namespace)
	}

	return &indexHandler{
		template: ai,
	}
}

func (ih *indexHandler) Info() *apiInfo {
	info := ih.template
	info.Env = os.Getenv("GO_ENV")
	return &info
}

func (ih *indexHandler) RegisterTo(r *hr.Router) {
	r.GET("/", ih.Handle)
	r.HEAD("/", ih.Handle)
}

func (ih *indexHandler) Handle(w http.ResponseWriter, _ *http.Request,
	_ hr.Params) {
	rw := responseWriter{w, ih.l}
	rw.WriteJSON(ih.Info())
}
