package controller

import (
	"github.com/mylxsw/glacier/infra"
	"github.com/mylxsw/glacier/web"
	"github.com/mylxsw/healthcheck/internal/config"
)

type InspectController struct {
	cc   infra.Resolver
	conf *config.Config
}

func NewInspectController(cc infra.Resolver, conf *config.Config) web.Controller {
	return &InspectController{cc: cc, conf: conf}
}

func (wel InspectController) Register(router web.Router) {
	router.Group("/inspect", func(router web.Router) {
		router.Any("/version", wel.Version)
	})
}

func (wel InspectController) Version(ctx web.Context) web.Response {
	return ctx.JSON(web.M{
		"version": wel.conf.Version,
		"git":     wel.conf.GitCommit,
	})
}
