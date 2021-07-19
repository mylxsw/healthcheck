package api

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/glacier/infra"
	"github.com/mylxsw/glacier/listener"
	"github.com/mylxsw/glacier/web"
	"github.com/mylxsw/healthcheck/api/controller"
)

type Provider struct{}

func (s Provider) Aggregates() []infra.Provider {
	return []infra.Provider{
		web.Provider(
			listener.FlagContext("listen"),
			web.SetRouteHandlerOption(s.routes),
			web.SetExceptionHandlerOption(s.exceptionHandler),
			web.SetIgnoreLastSlashOption(true),
		),
	}
}

func (s Provider) Register(app infra.Binder) {}
func (s Provider) Boot(app infra.Resolver)   {}

func (s Provider) exceptionHandler(ctx web.Context, err interface{}) web.Response {
	log.Errorf("error: %v, call stack: %s", err, debug.Stack())
	return ctx.JSONWithCode(web.M{
		"error": fmt.Sprintf("%v", err),
	}, http.StatusInternalServerError)
}

func (s Provider) routes(cc infra.Resolver, router web.Router, mw web.RequestMiddleware) {
	mws := make([]web.HandlerDecorator, 0)
	mws = append(mws,
		mw.AccessLog(log.Module("api")),
		mw.CORS("*"),
	)

	router.WithMiddleware(mws...).Controllers(
		"/api",
		controller.NewAlertController(cc),
		controller.NewHealthcheckController(cc),
	)
}
