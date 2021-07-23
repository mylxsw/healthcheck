package healthcheck

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/mylxsw/adanos-alert/pkg/connector"
	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/go-utils/failover"
)

// AlertType 告警类型
type AlertType string

const (
	// AlertTypeAdanos 基于 Adanos 的告警通知
	AlertTypeAdanos AlertType = "adanos"
	AlertTypeHTTP   AlertType = "http"
)

// AlertConfig 告警配置
type AlertConfig struct {
	Type        AlertType    `yaml:"type" json:"type"`
	AdanosAddr  string       `yaml:"adanos_addr" json:"adanos_addr"`
	AdanosToken string       `yaml:"adanos_token" json:"adanos_token"`
	HTTPAddr    string       `yaml:"http_addr" json:"http_addr"`
	HTTPHeaders []HTTPHeader `yaml:"http_header" json:"http_header"`
	Timeout     int64        `yaml:"timeout" json:"timeout"`
}

func (ac AlertConfig) SendEvent(ctx context.Context, status string, evt Event) error {
	retryer := failover.Retry(func(retryTimes int) error {
		switch ac.Type {
		case AlertTypeAdanos:
			return ac.sendAdanosAlert(ctx, status, evt)
		case AlertTypeHTTP:
			return ac.sendHTTPAlert(ctx, status, evt)
		}

		return fmt.Errorf("not support such alert type: %s", ac.Type)
	}, 3)

	_, err := retryer.Run()
	return err
}

func (ac AlertConfig) sendAdanosAlert(ctx context.Context, status string, evt Event) error {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(ac.Timeout)*time.Second)
	defer cancel()

	adanosEvt := connector.NewEvent(evt.Healthcheck.String()).
		WithOrigin("healthcheck").
		WithTags(evt.Healthcheck.Tags...).
		WithMeta("last_alert_time", evt.LastAlertTime.Format(time.RFC3339)).
		WithMeta("alert_times", evt.AlertTimes).
		WithMeta("last_failure", evt.LastFailure).
		WithMeta("last_success_time", evt.LastSuccessTime.Format(time.RFC3339)).
		WithMeta("healthcheck_id", evt.Healthcheck.ID).
		WithMeta("healthcheck_name", evt.Healthcheck.Name).
		WithMeta("status", status)

	conn := connector.NewConnector(ac.AdanosToken, ac.AdanosAddr)
	return conn.Send(ctx, adanosEvt)
}

func (ac AlertConfig) sendHTTPAlert(ctx context.Context, status string, evt Event) error {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(ac.Timeout)*time.Second)
	defer cancel()

	data, _ := json.Marshal(map[string]interface{}{
		"event":  evt,
		"status": status,
	})
	req, err := http.NewRequestWithContext(ctx, "POST", ac.HTTPAddr, bytes.NewReader(data))
	if err != nil {
		return err
	}

	for _, header := range ac.HTTPHeaders {
		req.Header.Add(header.Key, header.Value)
	}

	client := &http.Client{}
	client.Timeout = time.Duration(ac.Timeout) * time.Second

	if log.DebugEnabled() {
		log.With(evt).Debugf("send http alert message")
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("invalid status_code from http endpoint: %d %s", resp.StatusCode, resp.Status)
	}

	return nil
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
	LastSuccessTime time.Time `json:"last_success_time"`
}
