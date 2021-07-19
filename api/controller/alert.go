package controller

import (
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
	})
}

func (ctl AlertController) Alerts(alertManager *alert.Manager) []alert.Alert {
	return alertManager.GetAlerts()
}
