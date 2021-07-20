package scheduler

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/mylxsw/asteria/log"
)

// Scheduler 心跳检测调度器
type Scheduler struct {
	Jobs []*Job
	lock sync.RWMutex
}

// NewScheduler create a new scheduler
func NewScheduler() *Scheduler {
	return &Scheduler{
		Jobs: make([]*Job, 0),
	}
}

// PrintStatus 打印状态
func (sc *Scheduler) PrintStatus() {
	sc.lock.RLock()
	defer sc.lock.RUnlock()

	for _, job := range sc.Jobs {
		log.With(job).Debugf("job status")
	}
}

// AllJobs return all jobs
func (sc *Scheduler) AllJobs() HealthcheckJobs {
	jobs := make(HealthcheckJobs, 0)

	sc.lock.RLock()
	for _, job := range sc.Jobs {
		jobs = append(jobs, job.HealthcheckJob)
	}
	sc.lock.RUnlock()

	sort.Sort(jobs)
	return jobs
}

// RemoveJob remove a job
func (sc *Scheduler) RemoveJob(id string) {
	sc.lock.Lock()
	defer sc.lock.Unlock()

	idx := -1
	for i, job := range sc.Jobs {
		if job.Healthcheck.ID == id {
			idx = i
			break
		}
	}

	if idx >= 0 {
		if sc.Jobs[idx].Healthcheck.Editable {
			sc.Jobs = append(sc.Jobs[:idx], sc.Jobs[idx+1:]...)
		}
	}
}

// AddJob 添加 job
func (sc *Scheduler) AddJob(job *Job) {
	if log.DebugEnabled() {
		log.With(job.Healthcheck).Debugf("add healthcheck job: %s", job.Healthcheck.ID)
	}

	sc.lock.Lock()
	defer sc.lock.Unlock()

	sc.Jobs = append(sc.Jobs, job)
}

// Run 调度器执行
func (sc *Scheduler) Run(ctx context.Context) <-chan interface{} {
	stopped := make(chan interface{})
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				func() {
					sc.lock.RLock()
					defer sc.lock.RUnlock()

					for _, job := range sc.Jobs {
						if job.Schedulable() {
							go job.Run(ctx)
						}
					}
				}()
			case <-ctx.Done():
				stopped <- struct{}{}
				return
			}
		}
	}()

	return stopped
}
