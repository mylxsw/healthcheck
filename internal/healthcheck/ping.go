package healthcheck

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-ping/ping"
	"github.com/mylxsw/go-utils/failover"
)

// CheckTypeICMP ICMP Ping 类型的心跳检测
type CheckTypeICMP struct {
	Host        string `yaml:"host" json:"host"`
	Count       int64  `yaml:"count" json:"count"`
	Timeout     int64  `yaml:"timeout" json:"timeout"`
	SuccessRule string `yaml:"success_rule" json:"success_rule"`
}

func (cth CheckTypeICMP) init(timeout int64) CheckTypeICMP {
	if cth.Count == 0 {
		cth.Count = 3
	}

	if cth.Timeout == 0 {
		cth.Timeout = timeout
	}

	if cth.SuccessRule == "" {
		cth.SuccessRule = "PacketLoss == 0"
	}

	return cth
}

// Check 执行心跳检测
func (cth CheckTypeICMP) Check(ctx context.Context, hb Healthcheck) error {
	retryer := failover.Retry(func(retryTimes int) error {
		pinger, err := ping.NewPinger(cth.Host)
		if err != nil {
			return err
		}

		pinger.Timeout = time.Duration(cth.Timeout) * time.Second
		pinger.Count = int(cth.Count)
		pinger.Interval = 10 * time.Microsecond

		if err := pinger.Run(); err != nil {
			return err
		}

		data := pinger.Statistics()
		success, err := successRuleCheck(cth.SuccessRule, data)
		if err != nil {
			return err
		}

		if success {
			return nil
		}

		rs, _ := json.Marshal(data)
		return fmt.Errorf("success check return negative: %s", string(rs))
	}, 3)

	_, err := retryer.Run()
	return err
}
