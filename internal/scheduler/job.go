package scheduler

import (
	"context"
	"fmt"
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

type HealthcheckJobs []HealthcheckJob

func (jobs HealthcheckJobs) Len() int {
	return len(jobs)
}

func (jobs HealthcheckJobs) Less(i, j int) bool {
	return (jobs[i].Healthcheck.Name + jobs[i].Healthcheck.ID) < (jobs[j].Healthcheck.Name + jobs[j].Healthcheck.ID)
}

func (jobs HealthcheckJobs) Swap(i, j int) {
	jobs[i], jobs[j] = jobs[j], jobs[i]
}

type HealthcheckJob struct {
	Healthcheck     healthcheck.Healthcheck `json:"healthcheck,omitempty"`
	LastActiveTime  time.Time               `json:"last_active_time,omitempty"`
	LastFinishTime  time.Time               `json:"last_finish_time,omitempty"`
	LastSuccessTime time.Time               `json:"last_success_time,omitempty"`
	LastFailure     string                  `json:"last_failure,omitempty"`
}

// NewJob create a new job
func NewJob(hc healthcheck.Healthcheck) *Job {
	return &Job{
		HealthcheckJob: HealthcheckJob{
			Healthcheck:     hc,
			LastSuccessTime: time.Now(),
			LastFailure:     fmt.Sprintf("no heartbeat received since %s", time.Now()),
		},
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

	return job.Healthcheck.Schedulable() && time.Now().After(job.LastActiveTime.Add(time.Duration(job.Healthcheck.CheckInterval)*time.Second))
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
		job.lock.Unlock()

		log.With(job).Errorf("handle %s health check [%s] failed: %v", job.Healthcheck.CheckType, job.Healthcheck.Name, err)
	} else {
		job.lock.Lock()
		job.LastSuccessTime = time.Now()
		job.lock.Unlock()
	}

}

// UpdateJobStatus update job status
func (job *Job) UpdateJobStatus() {
	job.lock.Lock()
	defer job.lock.Unlock()

	job.LastActiveTime = time.Now()
	job.LastFinishTime = job.LastActiveTime
	job.LastSuccessTime = job.LastActiveTime

	job.LastFailure = fmt.Sprintf("no heartbeat received since %s", job.LastActiveTime)
}
