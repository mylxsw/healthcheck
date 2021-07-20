package healthcheck

import (
	"context"
	"fmt"
	"time"

	"github.com/mylxsw/adanos-alert/pkg/connector"
	"github.com/mylxsw/go-utils/failover"
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
	retryer := failover.Retry(func(retryTimes int) error {
		switch ac.Type {
		case AlertTypeAdanos:
			adanosEvt := connector.NewEvent(evt.Healthcheck.String()).
				WithOrigin("healthcheck").
				WithTags(evt.Healthcheck.Tags...).
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
	}, 3)

	_, err := retryer.Run()
	return err
}

type Event struct {
	Healthcheck Healthcheck `json:"healthcheck"`
	// LastAliveTime 最后一次该心跳检测活跃的时间，该字段用于检测一个 Alert 是否已经失效了
	LastAliveTime time.Time `json:"last_alive_time"`

	// lastAlertTime 最后一次告警时间
	LastAlertTime time.Time `json:"last_alert_time"`
	// alertTimes 告警次数，从最后一次心跳丢失开始
	AlertTimes int64 `json:"alert_times"`

	LastFailure     string    `json:"last_failure"`
	LastFailureTime time.Time `json:"last_failure_time"`
	LastSuccessTime time.Time `json:"last_success_time"`
}
