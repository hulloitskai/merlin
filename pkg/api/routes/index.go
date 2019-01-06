package routes

import (
	"fmt"
	"net/http"
	"os"

	"go.uber.org/zap"

	hr "github.com/julienschmidt/httprouter"
	"github.com/stevenxie/merlin/internal/info"
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
	info := apiInfo{
		ID:      fmt.Sprintf("%s-api-%s", info.Namespace, info.Version),
		Version: info.Version,
	}
	return &indexHandler{
		template: info,
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
