package controller

import (
	"net/http"

	"github.com/mylxsw/glacier/infra"
	"github.com/mylxsw/glacier/web"
	"github.com/mylxsw/healthcheck/internal/scheduler"
)

type PushController struct {
	cc infra.Resolver
}

func NewPushController(cc infra.Resolver) web.Controller {
	return &PushController{cc: cc}
}

func (ctl PushController) Register(router web.Router) {
	router.Group("/push", func(router web.Router) {
		router.Any("/{id}", ctl.Push)
		router.Any("/{id}/", ctl.Push)
	})
}

func (ctl PushController) Push(webCtx web.Context, sche *scheduler.Scheduler) web.Response {
	id := webCtx.PathVar("id")
	if err := sche.UpdateJobStatus(id); err != nil {
		return webCtx.JSONError(err.Error(), http.StatusBadRequest)
	}

	return webCtx.JSON(web.M{"id": id})
}
