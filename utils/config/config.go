package config

import (
	"github.com/spf13/viper"
)

//AppName global variable
const AppName = "uv-server"

//Manager type
type Manager struct{ *viper.Viper }

var config *viper.Viper

//Config get global configs
func Config() *Manager { return &Manager{config} }

func init() {
	config = viper.New()
	config.AutomaticEnv()
}
