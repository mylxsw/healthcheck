package healthcheck

import (
	"context"
	"fmt"
	"time"

	"github.com/mylxsw/adanos-alert/pkg/connector"
)

// AlertType 告警类型
type AlertType string

const (
	// AlertTypeAdanos 基于 Adanos 的告警通知
	AlertTypeAdanos AlertType = "adanos"
)

// AlertConfig 告警配置
type AlertConfig struct {
	Type        AlertType `yaml:"type" json:"type"`
	AdanosAddr  string    `yaml:"adanos_addr" json:"adanos_addr"`
	AdanosToken string    `yaml:"adanos_token" json:"adanos_token"`
}

func (ac AlertConfig) SendEvent(ctx context.Context, status string, evt Event) error {
	switch ac.Type {
	case AlertTypeAdanos:
		adanosEvt := connector.NewEvent(evt.Healthcheck.String()).
			WithOrigin("healthcheck").
			WithMeta("last_alert_time", evt.LastAlertTime.Format(time.RFC3339)).
			WithMeta("alert_times", evt.AlertTimes).
			WithMeta("last_failure", evt.LastFailure).
			WithMeta("last_failure_time", evt.LastFailureTime.Format(time.RFC3339)).
			WithMeta("last_success_time", evt.LastSuccessTime.Format(time.RFC3339)).
			WithMeta("healthcheck_id", evt.Healthcheck.ID).
			WithMeta("healthcheck_name", evt.Healthcheck.Name).
			WithMeta("status", status)

		conn := connector.NewConnector(ac.AdanosToken, ac.AdanosAddr)
		return conn.Send(ctx, adanosEvt)
	}

	return fmt.Errorf("not support such alert type: %s", ac.Type)
}

type Event struct {
	Healthcheck Healthcheck
	// LastAliveTime 最后一次该心跳检测活跃的时间，该字段用于检测一个 Alert 是否已经失效了
	LastAliveTime time.Time

	// lastAlertTime 最后一次告警时间
	LastAlertTime time.Time
	// alertTimes 告警次数，从最后一次心跳丢失开始
	AlertTimes int64

	LastFailure     string
	LastFailureTime time.Time
	LastSuccessTime time.Time
}
