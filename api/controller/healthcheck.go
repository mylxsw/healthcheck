package controller

import (
	"github.com/mylxsw/glacier/infra"
	"github.com/mylxsw/glacier/web"
	"github.com/mylxsw/healthcheck/internal/scheduler"
)

type HealthcheckController struct {
	cc infra.Resolver
}

func NewHealthcheckController(cc infra.Resolver) web.Controller {
	return &HealthcheckController{cc: cc}
}

func (ctl HealthcheckController) Register(router web.Router) {
	router.Group("/healthchecks", func(router web.Router) {
		router.Get("/", ctl.HealthChecks)
		router.Delete("/{id}", ctl.DeleteHealthCheck)
	})
}

func (ctl HealthcheckController) HealthChecks(sche *scheduler.Scheduler) []scheduler.HealthcheckJob {
	return sche.AllJobs()
}

func (ctl HealthcheckController) DeleteHealthCheck(webCtx web.Context, sche *scheduler.Scheduler) web.Response {
	id := webCtx.PathVar("id")
	sche.RemoveJob(id)

	return webCtx.JSON(web.M{})
}
