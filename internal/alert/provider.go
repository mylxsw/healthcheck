package alert

import (
	"context"

	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/glacier/infra"
	"github.com/mylxsw/healthcheck/internal/healthcheck"
	"github.com/mylxsw/healthcheck/internal/scheduler"
)

type Provider struct{}

func (s Provider) Register(cc infra.Binder) {
	cc.MustSingletonOverride(func(globalConf *healthcheck.GlobalConfig, sche *scheduler.Scheduler) *Manager {
		return NewManager(globalConf, sche, 1024)
	})
}
func (s Provider) Boot(cc infra.Resolver) {}

func (s Provider) Daemon(ctx context.Context, app infra.Resolver) {
	app.MustResolve(func(m *Manager) {
		<-m.Run(ctx)
		if log.DebugEnabled() {
			log.Debugf("alert manager has been stopped")
		}
	})
}
