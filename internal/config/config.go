package config

import "github.com/mylxsw/glacier/infra"

type Config struct {
	Listen    string
	Version   string
	GitCommit string
}

// Get return config object from container
func Get(cc infra.Resolver) *Config {
	return cc.MustGet(&Config{}).(*Config)
}
