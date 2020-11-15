package config

import (
	"encoding/json"

	"github.com/mylxsw/container"
)

type Config struct {
	Listen  string `json:"listen"`
	Debug   bool   `json:"debug"`
}

func (conf *Config) Serialize() string {
	rs, _ := json.Marshal(conf)
	return string(rs)
}

// Get return config object from container
func Get(cc container.Container) *Config {
	return cc.MustGet(&Config{}).(*Config)
}
