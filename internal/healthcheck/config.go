package healthcheck

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"gopkg.in/yaml.v2"
)

// GlboalConfig is a global configuration object
type GlobalConfig struct {
	Version       string        `yaml:"version" json:"version"`
	Healthchecks  []Healthcheck `yaml:"healthchecks" json:"healthchecks"`
	Discoveries   []Discovery   `yaml:"discoveries" json:"discoveries"`
	WorkerNum     int           `yaml:"worker_num" json:"worker_num"`
	CheckInterval int64         `yaml:"check_interval" json:"check_interval"`
	LossThreshold int64         `yaml:"loss_threshold" json:"loss_threshold"`
	HTTPTimeout   int64         `yaml:"http_timeout" json:"http_timeout"`
	PINGTimeout   int64         `yaml:"ping_timeout" json:"ping_timeout"`
	Alerts        []AlertConfig `yaml:"alerts" json:"alerts"`
}

func (gc *GlobalConfig) init() error {
	if gc.Version == "" {
		gc.Version = "1"
	}

	if gc.WorkerNum == 0 {
		gc.WorkerNum = 3
	}

	if gc.CheckInterval == 0 {
		gc.CheckInterval = 60
	}

	if gc.LossThreshold == 0 {
		gc.LossThreshold = gc.CheckInterval * 2
	}

	if gc.HTTPTimeout == 0 {
		gc.HTTPTimeout = 60
	}

	if gc.PINGTimeout == 0 {
		gc.PINGTimeout = 1
	}

	if gc.Healthchecks == nil {
		gc.Healthchecks = make([]Healthcheck, 0)
	}

	for i, hb := range gc.Healthchecks {
		hb.Editable = false
		if hb.Tags == nil {
			gc.Healthchecks[i].Tags = make([]string, 0)
		}

		if hb.CheckType == "" {
			return fmt.Errorf("invalid check_type")
		}

		if hb.CheckInterval == 0 {
			gc.Healthchecks[i].CheckInterval = gc.CheckInterval
		}

		if hb.LossThreshold == 0 {
			gc.Healthchecks[i].LossThreshold = gc.LossThreshold
		}

		gc.Healthchecks[i].ID = fmt.Sprintf("check-%s", hb.CheckType)
		if hb.Name == "" {
			gc.Healthchecks[i].Name = gc.Healthchecks[i].ID
		}

		switch hb.CheckType {
		case HTTP:
			gc.Healthchecks[i].HTTP = hb.HTTP.init(gc.HTTPTimeout)
		case PING:
			gc.Healthchecks[i].PING = hb.PING.init(gc.PINGTimeout)
		default:
		}
	}

	for i, alt := range gc.Alerts {
		if alt.Timeout == 0 {
			alt.Timeout = 30
		}

		if alt.SilentPeriod == "" {
			alt.SilentPeriod = "1m"
		}

		_, err := time.ParseDuration(alt.SilentPeriod)
		if err != nil {
			return fmt.Errorf("invalid alert.silent_period: %v", err)
		}

		if alt.HTTPHeaders == nil {
			alt.HTTPHeaders = make([]HTTPHeader, 0)
		}

		gc.Alerts[i] = alt
	}

	return nil
}

// ToYAML return GlobalConfig as yaml
func (gc *GlobalConfig) ToYAML() string {
	data, _ := yaml.Marshal(gc)
	return string(data)
}

// CheckType 健康检查类型
type CheckType string

const (
	// HTTP http 类型的健康检查
	HTTP CheckType = "http"
	PING CheckType = "ping"
	PUSH CheckType = "push"
)

// Healthcheck 健康检查对象
type Healthcheck struct {
	ID            string        `yaml:"-" json:"id,omitempty"`
	Editable      bool          `yaml:"-" json:"editable,omitempty"`
	Name          string        `yaml:"name" json:"name,omitempty"`
	Tags          []string      `yaml:"tags" json:"tags,omitempty"`
	CheckInterval int64         `yaml:"check_interval" json:"check_interval,omitempty"`
	LossThreshold int64         `yaml:"loss_threshold" json:"loss_threshold,omitempty"`
	CheckType     CheckType     `yaml:"check_type" json:"check_type,omitempty"`
	HTTP          CheckTypeHTTP `yaml:"http" json:"http,omitempty"`
	PING          CheckTypeICMP `yaml:"ping" json:"ping,omitempty"`
}

// String convert healthcheck to string
func (hb Healthcheck) String() string {
	data, _ := json.Marshal(hb)
	return string(data)
}

// Schedulable return whether the Healthcheck is schedulable
func (hb Healthcheck) Schedulable() bool {
	return hb.CheckType != PUSH
}

// Check 发起健康检查
func (hb Healthcheck) Check(ctx context.Context) error {
	switch hb.CheckType {
	case HTTP:
		return hb.HTTP.Check(ctx, hb)
	case PING:
		return hb.PING.Check(ctx, hb)
	}

	return nil
}

// ParseYamlConfig parse config from yaml
func ParseYamlConfig(data []byte) (*GlobalConfig, error) {
	var globalConf GlobalConfig
	if err := yaml.Unmarshal(data, &globalConf); err != nil {
		return nil, err
	}

	if err := globalConf.init(); err != nil {
		return nil, err
	}

	for i, dis := range globalConf.Discoveries {
		globalConf.Discoveries[i] = dis.init(globalConf)
	}

	return &globalConf, nil
}
