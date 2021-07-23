package healthcheck

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/go-utils/failover"
	"github.com/mylxsw/pattern"
)

// HTTPHeader http 请求头
type HTTPHeader struct {
	Key   string `yaml:"key" json:"key"`
	Value string `yaml:"value" json:"value"`
}

// CheckTypeHTTP HTTP 类型的心跳检测
type CheckTypeHTTP struct {
	Method      string       `yaml:"method" json:"method"`
	URL         string       `yaml:"url" json:"url"`
	Headers     []HTTPHeader `yaml:"headers" json:"-"`
	Body        string       `yaml:"body" json:"body"`
	Timeout     int64        `yaml:"timeout" json:"timeout"`
	SuccessRule string       `yaml:"success_rule" json:"success_rule"`
}

func (cth CheckTypeHTTP) init(timeout int64) CheckTypeHTTP {
	if cth.Method == "" {
		cth.Method = "GET"
	}

	if cth.Headers == nil {
		cth.Headers = make([]HTTPHeader, 0)
	}

	if cth.Timeout == 0 {
		cth.Timeout = timeout
	}

	if cth.SuccessRule == "" {
		cth.SuccessRule = "StatusCode >= 200 and StatusCode < 400"
	}

	return cth
}

// Check 执行心跳检测
func (cth CheckTypeHTTP) Check(ctx context.Context, hb Healthcheck) error {
	retryer := failover.Retry(func(retryTimes int) error {
		ctx, cancel := context.WithTimeout(ctx, time.Duration(cth.Timeout)*time.Second)
		defer cancel()

		req, err := http.NewRequestWithContext(ctx, cth.Method, cth.URL, strings.NewReader(cth.Body))
		if err != nil {
			return err
		}

		for _, header := range cth.Headers {
			req.Header.Add(header.Key, header.Value)
		}

		client := &http.Client{}
		client.Timeout = time.Duration(cth.Timeout) * time.Second

		if log.DebugEnabled() {
			log.With(hb).Debugf("send http health check")
		}

		resp, err := client.Do(req)
		if err != nil {
			return err
		}

		defer resp.Body.Close()
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		checkData := SuccessRuleCheckData{
			StatusCode: resp.StatusCode,
			Status:     resp.Status,
			Body:       string(respBody),
		}

		if log.DebugEnabled() {
			log.WithFields(log.Fields{
				"healthcheck": hb,
				"response":    checkData,
			}).Debugf("http healthcheck response")
		}

		success, err := successRuleCheck(cth.SuccessRule, checkData)
		if err != nil {
			return err
		}

		if success {
			return nil
		}

		return fmt.Errorf("success check return negative: %s", checkData.String())
	}, 3)

	_, err := retryer.Run()
	return err
}

// SuccessRuleCheckData 请求结果判定
type SuccessRuleCheckData struct {
	pattern.Helpers
	StatusCode int
	Status     string
	Body       string
}

// String 请求结果判定对象字符串表示
func (srcd SuccessRuleCheckData) String() string {
	data, _ := json.Marshal(srcd)
	return string(data)
}

// successRuleCheck 请求结果判定
func successRuleCheck(rule string, data interface{}) (bool, error) {
	matcher, err := pattern.NewMatcher(rule, data)
	if err != nil {
		return false, err
	}

	return matcher.Match(data)
}
