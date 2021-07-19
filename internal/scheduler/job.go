package scheduler

import (
	"context"
	"sync"
	"time"

	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/healthcheck/internal/healthcheck"
)

// Job 心跳检测任务对象
type Job struct {
	HealthcheckJob
	lock sync.RWMutex
}

type HealthcheckJob struct {
	Healthcheck     healthcheck.Healthcheck `json:"healthcheck"`
	LastActiveTime  time.Time               `json:"last_active_time"`
	LastFinishTime  time.Time               `json:"last_finish_time"`
	LastSuccessTime time.Time               `json:"last_success_time"`
	LastFailure     string                  `json:"last_failure"`
	LastFailureTime time.Time               `json:"last_failure_time"`
}

// NewJob create a new job
func NewJob(hc healthcheck.Healthcheck) *Job {
	return &Job{
		HealthcheckJob: HealthcheckJob{Healthcheck: hc, LastSuccessTime: time.Now()},
	}
}

// Failed 返回当前心跳检测是否失败（心跳丢失时间大于 LossThreshold）
func (job HealthcheckJob) Failed() bool {
	return time.Now().After(job.LastSuccessTime.Add(time.Duration(job.Healthcheck.LossThreshold) * time.Second))
}

// Schedulable 返回心跳检测对象是否符合检测条件
func (job *Job) Schedulable() bool {
	job.lock.RLock()
	defer job.lock.RUnlock()

	return time.Now().After(job.LastActiveTime.Add(time.Duration(job.Healthcheck.CheckInterval) * time.Second))
}

// Run 执行任务
func (job *Job) Run(ctx context.Context) {
	if !job.Schedulable() {
		return
	}

	job.lock.Lock()
	job.LastActiveTime = time.Now()
	job.lock.Unlock()

	defer func() {
		job.lock.Lock()
		job.LastFinishTime = time.Now()
		job.lock.Unlock()
	}()

	if err := job.Healthcheck.Check(ctx); err != nil {
		job.lock.Lock()
		job.LastFailure = err.Error()
		job.LastFailureTime = time.Now()
		job.lock.Unlock()

		log.With(job).Errorf("handle %s health check [%s] failed: %v", job.Healthcheck.CheckType, job.Healthcheck.Name, err)
	} else {
		job.lock.Lock()
		job.LastSuccessTime = time.Now()
		job.lock.Unlock()
	}

}
