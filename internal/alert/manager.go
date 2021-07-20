package alert

import (
	"context"
	"sync"
	"time"

	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/healthcheck/internal/healthcheck"
	"github.com/mylxsw/healthcheck/internal/scheduler"
)

type EventType string

const (
	EventTypeSuccess EventType = "success"
	EventTypeFail    EventType = "fail"
)

type Event struct {
	Type  EventType `json:"type"`
	Alert Alert     `json:"alert"`
}

type Manager struct {
	alerts map[string]*Alert
	sche   *scheduler.Scheduler
	queue  chan Event
	conf   *healthcheck.GlobalConfig

	lock sync.RWMutex
}

func NewManager(globalConf *healthcheck.GlobalConfig, sche *scheduler.Scheduler, queueSize int64) *Manager {
	return &Manager{conf: globalConf, alerts: make(map[string]*Alert), sche: sche, queue: make(chan Event, queueSize)}
}

func (m *Manager) GetAlerts() []Alert {
	m.lock.RLock()
	defer m.lock.RUnlock()

	alerts := make([]Alert, 0)
	for _, al := range m.alerts {
		alerts = append(alerts, *al)
	}

	return alerts
}

func (m *Manager) Run(ctx context.Context) <-chan interface{} {
	stopped := make(chan interface{})
	safeClose := make(chan interface{})

	go func() {
		for evt := range m.queue {
			if evt.Type == EventTypeFail {
				log.With(evt.Alert).Errorf("healthcheck for %s failed", evt.Alert.Healthcheck.ID)
			} else {
				log.With(evt.Alert).Infof("healthcheck for %s succeed", evt.Alert.Healthcheck.ID)
			}

			for i, alt := range m.conf.Alerts {
				if err := alt.SendEvent(ctx, string(evt.Type), evt.Alert.Event); err != nil {
					log.With(evt).Errorf("send event to alert channel %s-%d failed: %v", alt.Type, i, err)
				}
			}
		}

		safeClose <- struct{}{}
	}()

	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				m.handleHealthCheck()
				m.collect()
			case <-ctx.Done():
				close(m.queue)
				<-safeClose
				close(safeClose)

				stopped <- struct{}{}
				return
			}
		}
	}()

	return stopped
}

func (m *Manager) handleHealthCheck() {
	for _, job := range m.sche.AllJobs() {
		m.lock.Lock()

		alert, ok := m.alerts[job.Healthcheck.ID]
		if !ok {
			alert = NewAlert(job.Healthcheck)
			m.alerts[job.Healthcheck.ID] = alert
		}

		alert.LastAliveTime = time.Now()
		if alert.IsFailed() && !job.Failed() {
			// 失败 -> 成功
			alert.MarkSucceed(job.LastSuccessTime)
			select {
			case m.queue <- Event{Type: EventTypeSuccess, Alert: *alert}:
			default:
			}
		} else if !alert.IsFailed() && job.Failed() {
			// 成功 -> 失败
			alert.MarkFailed(job.LastFailureTime, job.LastFailure)
			select {
			case m.queue <- Event{Type: EventTypeFail, Alert: *alert}:
			default:
			}
		}

		m.lock.Unlock()
	}
}

// collect 删除已经不存在的健康检查
func (m *Manager) collect() {
	m.lock.Lock()
	defer m.lock.Unlock()

	expired := make([]string, 0)
	for k, v := range m.alerts {
		// 删除持续 60s 没有被触发过的告警配置
		if time.Since(v.LastAliveTime).Seconds() > 60 {
			expired = append(expired, k)
		}
	}

	for _, k := range expired {
		delete(m.alerts, k)
	}
}

type Alert struct {
	healthcheck.Event
}

func NewAlert(hb healthcheck.Healthcheck) *Alert {
	return &Alert{Event: healthcheck.Event{Healthcheck: hb}}
}

func (alert *Alert) MarkSucceed(successTime time.Time) {
	alert.AlertTimes = 0
	alert.LastSuccessTime = successTime

}

func (alert *Alert) IsFailed() bool {
	return alert.AlertTimes > 0
}

func (alert *Alert) MarkFailed(failureTime time.Time, failureReason string) {
	alert.LastFailure = failureReason
	alert.LastFailureTime = failureTime

	if alert.AlertTimes > 0 {
		return
	}

	alert.AlertTimes = 1
	alert.LastAlertTime = time.Now()
}
