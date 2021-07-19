package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/mylxsw/asteria/formatter"
	"github.com/mylxsw/asteria/level"
	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/container"
	"github.com/mylxsw/glacier/infra"
	"github.com/mylxsw/glacier/starter/application"
	"github.com/mylxsw/healthcheck/api"
	"github.com/mylxsw/healthcheck/internal/alert"
	"github.com/mylxsw/healthcheck/internal/config"
	"github.com/mylxsw/healthcheck/internal/healthcheck"
	"github.com/mylxsw/healthcheck/internal/scheduler"
	"github.com/urfave/cli"
	"github.com/urfave/cli/altsrc"
)

var Version = "1.0"
var GitCommit = "5dbef13fb456f51a5d29464d"

func main() {
	app := application.Create(fmt.Sprintf("%s %s", Version, GitCommit))
	app.AddFlags(altsrc.NewStringFlag(cli.StringFlag{
		Name:  "listen",
		Usage: "服务监听地址",
		Value: "127.0.0.1:10101",
	}))
	app.AddFlags(altsrc.NewStringFlag(cli.StringFlag{
		Name:  "healthcheck",
		Usage: "健康检查配置文件路径",
		Value: "healthchecks.yaml",
	}))
	app.AddFlags(altsrc.NewBoolFlag(cli.BoolFlag{
		Name:  "debug",
		Usage: "是否使用调试模式",
	}))

	app.BeforeServerStart(func(cc container.Container) error {
		return cc.Resolve(func(c infra.FlagContext) {
			if !c.Bool("debug") {
				log.All().LogLevel(level.Info)
				log.All().LogFormatter(formatter.NewJSONFormatter())
			}
		})
	})

	app.Singleton(func(c infra.FlagContext) *healthcheck.GlobalConfig {
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

	app.Singleton(func(c infra.FlagContext) *config.Config {
		return &config.Config{
			Version:   Version,
			GitCommit: GitCommit,
			Listen:    c.String("listen"),
		}
	})

	app.Provider(scheduler.Provider{}, alert.Provider{}, api.Provider{})

	if err := app.Run(os.Args); err != nil {
		log.Errorf("exit with error: %s", err)
	}
}
