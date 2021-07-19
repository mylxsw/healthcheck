package healthcheck

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"text/template"

	"github.com/hashicorp/consul/api"
	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/pattern"
)

// DiscoveryType 服务发现类型
type DiscoveryType string

const (
	// DiscoveryTypeConsul 基于 Consul 的服务发现
	DiscoveryTypeConsul DiscoveryType = "consul"
)

// Discovery 服务发现
type Discovery struct {
	Type         DiscoveryType `yaml:"type" json:"type"`
	ConsulScheme string        `yaml:"consul_scheme" json:"consul_scheme"`
	ConsulAddr   string        `yaml:"consul_addr" json:"consul_addr"`
	ConsulToken  string        `yaml:"consul_token" json:"consul_token"`
	ConsulDC     string        `yaml:"consul_dc" json:"consul_dc"`
	Filter       string        `yaml:"filter" json:"filter"`
	Template     Healthcheck   `yaml:"template" json:"template"`
}

func (dis Discovery) init(conf GlobalConfig) Discovery {
	if dis.ConsulScheme == "" {
		dis.ConsulScheme = "http"
	}

	if dis.ConsulAddr == "" {
		dis.ConsulAddr = "127.0.0.1:8500"
	}

	if dis.Filter == "" {
		dis.Filter = "true"
	}

	if dis.Template.CheckInterval == 0 {
		dis.Template.CheckInterval = conf.CheckInterval
	}

	if dis.Template.LossThreshold == 0 {
		dis.Template.LossThreshold = conf.LossThreshold
	}

	if dis.Template.HTTP.Method == "" {
		dis.Template.HTTP.Method = "GET"
	}

	if dis.Template.HTTP.SuccessRule == "" {
		dis.Template.HTTP.SuccessRule = "StatusCode >= 200 and StatusCode < 400"
	}

	if dis.Template.HTTP.Timeout == 0 {
		dis.Template.HTTP.Timeout = conf.HTTPTimeout
	}

	return dis
}

func (dis Discovery) LoadHealthchecks(ctx context.Context, conf *GlobalConfig) ([]Healthcheck, error) {
	switch dis.Type {
	case DiscoveryTypeConsul:
		return dis.handleLoadHealthchecksFromConsul(ctx, conf)
	}

	return []Healthcheck{}, nil
}

func (dis Discovery) handleLoadHealthchecksFromConsul(ctx context.Context, conf *GlobalConfig) ([]Healthcheck, error) {
	client, err := api.NewClient(&api.Config{
		Scheme:     dis.ConsulScheme,
		Address:    dis.ConsulAddr,
		Token:      dis.ConsulToken,
		Datacenter: dis.ConsulDC,
	})
	if err != nil {
		return nil, err
	}

	services, _, err := client.Catalog().Services(nil)
	if err != nil {
		return nil, err
	}

	healthchecks := make([]Healthcheck, 0)

	for name, tags := range services {
		instances, _, err := client.Catalog().Service(name, "", nil)
		if err != nil {
			return nil, err
		}

		for _, ins := range instances {
			// filter
			srvObj := ConsulService{
				Tags:           tags,
				CatalogService: ins,
			}
			matched, err := srvObj.matched(dis.Filter)
			if err != nil {
				log.WithFields(log.Fields{"service": srvObj, "discovery": dis}).Errorf("eval filter rule failed: %v", err)
				continue
			}

			if !matched {
				continue
			}

			// create healthcheck
			hc := Healthcheck{
				ID:            fmt.Sprintf("discovery-consul-%s", ins.ServiceID),
				Editable:      true,
				CheckInterval: dis.Template.CheckInterval,
				LossThreshold: dis.Template.LossThreshold,
				CheckType:     dis.Template.CheckType,
				HTTP:          dis.Template.HTTP,
			}

			switch dis.Template.CheckType {
			case HTTP:
				hc.HTTP.Timeout = dis.Template.HTTP.Timeout

				if dis.Template.Name == "" {
					dis.Template.Name = "tmpl: {{ .ServiceName }}:{{ .ServiceAddress }}:{{ .ServicePort }}"
				}

				hc.Name, err = srvObj.parseTemplate(dis.Template.Name)
				if err != nil {
					log.With(log.Fields{"service": srvObj, "discovery": dis}).Errorf("parse template for name failed: %v", err)
				}

				hc.HTTP.Method, err = srvObj.parseTemplate(dis.Template.HTTP.Method)
				if err != nil {
					log.With(log.Fields{"service": srvObj, "discovery": dis}).Errorf("parse template for http.method failed: %v", err)
				}

				if dis.Template.HTTP.URL == "" {
					dis.Template.HTTP.URL = "http://{{ .ServiceAddress }}:{{ .ServicePort }}/health"
				}

				hc.HTTP.URL, err = srvObj.parseTemplate(dis.Template.HTTP.URL)
				if err != nil {
					log.With(log.Fields{"service": srvObj, "discovery": dis}).Errorf("parse template for http.url failed: %v", err)
				}

				hc.HTTP.Body, err = srvObj.parseTemplate(dis.Template.HTTP.Body)
				if err != nil {
					log.With(log.Fields{"service": srvObj, "discovery": dis}).Errorf("parse template for http.body failed: %v", err)
				}

				hc.HTTP.SuccessRule, err = srvObj.parseTemplate(dis.Template.HTTP.SuccessRule)
				if err != nil {
					log.With(log.Fields{"service": srvObj, "discovery": dis}).Errorf("parse template for http.success_rule failed: %v", err)
				}

				hc.HTTP.Headers = make([]HTTPHeader, 0)
				for _, header := range dis.Template.HTTP.Headers {
					newHeader := HTTPHeader{}
					newHeader.Key, err = srvObj.parseTemplate(header.Key)
					if err != nil {
						log.With(log.Fields{"service": srvObj, "discovery": dis}).Errorf("parse template for http.headers.%s failed: %v", header.Key, err)
					}

					newHeader.Value, err = srvObj.parseTemplate(header.Value)
					if err != nil {
						log.With(log.Fields{"service": srvObj, "discovery": dis}).Errorf("parse template for http.headers.%s=%s failed: %v", header.Key, header.Value, err)
					}

					hc.HTTP.Headers = append(hc.HTTP.Headers, newHeader)
				}

				hc.HTTP = hc.HTTP.init(conf.HTTPTimeout)
			}

			healthchecks = append(healthchecks, hc)
		}
	}

	return healthchecks, nil
}

type ConsulService struct {
	Tags []string
	*api.CatalogService
}

func (filter ConsulService) parseTemplate(tmp string) (string, error) {
	if !strings.HasPrefix(tmp, "tmpl:") {
		return tmp, nil
	}

	content := strings.TrimSpace(strings.TrimPrefix(tmp, "tmpl:"))

	parser, err := template.New("").Parse(content)
	if err != nil {
		return "", err
	}

	var buffer bytes.Buffer
	if err := parser.Execute(&buffer, filter); err != nil {
		return "", err
	}

	return buffer.String(), nil
}

func (filter ConsulService) matched(expr string) (bool, error) {
	matcher, err := pattern.NewMatcher(expr, filter)
	if err != nil {
		return false, err
	}

	return matcher.Match(filter)
}
