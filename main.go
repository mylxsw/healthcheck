package main

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/mylxsw/asteria/formatter"
	"github.com/mylxsw/asteria/level"
	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/glacier/infra"
	"github.com/mylxsw/glacier/starter/app"
	"github.com/mylxsw/healthcheck/api"
	"github.com/mylxsw/healthcheck/internal/alert"
	"github.com/mylxsw/healthcheck/internal/config"
	"github.com/mylxsw/healthcheck/internal/healthcheck"
	"github.com/mylxsw/healthcheck/internal/scheduler"
)

var Version = "1.0"
var GitCommit = "5dbef13fb456f51a5d29464d"

func main() {
	ins := app.Create(fmt.Sprintf("%s %s", Version, GitCommit), 3).WithShutdownTimeoutFlag(3 * time.Second)

	ins.AddStringFlag("listen", "127.0.0.1:10101", "服务监听地址")
	ins.AddStringFlag("healthcheck", "healthchecks.yaml", "健康检查配置文件路径")
	ins.AddBoolFlag("debug", "是否使用调试模式")

	ins.BeforeServerStop(func(resolver infra.Resolver) error {
		return resolver.Resolve(func(c infra.FlagContext) {
			if !c.Bool("debug") {
				log.All().LogLevel(level.Info)
				log.All().LogFormatter(formatter.NewJSONFormatter())
			}
		})
	})

	ins.Singleton(func(c infra.FlagContext) *healthcheck.GlobalConfig {
		confData, err := ioutil.ReadFile(c.String("healthcheck"))
		if err != nil {
			panic(fmt.Errorf("read config file from %s failed: %v", c.String("healthcheck"), err))
		}

		// healthcheck
		globalConf, err := healthcheck.ParseYamlConfig(confData)
		if err != nil {
			panic(fmt.Errorf("parse globalConfig failed: %v", err))
		}

		return globalConf
	})

	ins.Singleton(func(c infra.FlagContext) *config.Config {
		return &config.Config{
			Version:   Version,
			GitCommit: GitCommit,
			Listen:    c.String("listen"),
		}
	})

	ins.Provider(scheduler.Provider{}, alert.Provider{}, api.Provider{})

	app.MustRun(ins)
}
