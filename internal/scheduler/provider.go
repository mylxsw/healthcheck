package scheduler

import (
	"context"
	"time"

	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/coll"
	"github.com/mylxsw/glacier/infra"
	"github.com/mylxsw/healthcheck/internal/healthcheck"
)

type Provider struct{}

func (s Provider) Register(cc infra.Binder) {
	cc.MustSingletonOverride(func(ctx context.Context, globalConf *healthcheck.GlobalConfig) *Scheduler {
		sche := NewScheduler()
		for _, hb := range globalConf.Healthchecks {
			sche.AddJob(NewJob(hb))
		}

		addDiscoveriedJobs := func() {
			var jobs map[string]HealthcheckJob
			_ = coll.MustNew(sche.AllJobs()).AsMap(func(job HealthcheckJob) string { return job.Healthcheck.ID }).All(&jobs)
			for _, dis := range globalConf.Discoveries {
				hs, err := dis.LoadHealthchecks(globalConf)
				if err != nil {
					log.With(dis).Errorf("load healthchecks from discovery failed: %v", err)
					continue
				}

				for _, h := range hs {
					if _, ok := jobs[h.ID]; !ok {
						sche.AddJob(NewJob(h))
					}
				}
			}
		}

		addDiscoveriedJobs()
		go func() {
			ticker := time.NewTicker(5 * 60 * time.Second)
			defer ticker.Stop()

			for {
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
					addDiscoveriedJobs()
				}
			}
		}()

		return sche
	})
}
func (s Provider) Boot(cc infra.Resolver) {}

func (s Provider) Daemon(ctx context.Context, app infra.Resolver) {
	app.MustResolve(func(sche *Scheduler) {
		<-sche.Run(ctx)
	})
}
