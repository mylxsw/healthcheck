package controller

import (
	"github.com/mylxsw/coll"
	"github.com/mylxsw/glacier/infra"
	"github.com/mylxsw/glacier/web"
	"github.com/mylxsw/healthcheck/internal/alert"
)

type AlertController struct {
	cc infra.Resolver
}

func NewAlertController(cc infra.Resolver) web.Controller {
	return &AlertController{cc: cc}
}

func (ctl AlertController) Register(router web.Router) {
	router.Group("/alerts", func(router web.Router) {
		router.Get("/", ctl.Alerts)
		router.Get("/failed-count/", ctl.FailedCount)
	})
}

func (ctl AlertController) Alerts(alertManager *alert.Manager) alert.Alerts {
	return alertManager.GetAlerts()
}

func (ctl AlertController) FailedCount(webCtx web.Context, alertManager *alert.Manager) web.Response {
	return webCtx.JSON(web.M{
		"count": coll.MustNew(alertManager.GetAlerts()).Filter(func(a alert.Alert) bool { return a.AlertTimes > 0 }).Size(),
	})
}
