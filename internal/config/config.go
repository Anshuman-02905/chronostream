package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Engine struct {
		ProducerVersion string
		InstanceID      string
		BufferSize      int
		Message         string
	}
	Dispatcher struct {
		MaxRetries  int
		BaseBackoff int
		MaxBackoff  int
	}
}

func (c *Config) Load() {
	viper.SetConfigName("config")
	viper.AddConfigPath("./config")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	c.Engine.ProducerVersion = viper.GetString("engine.producer_version")
	c.Engine.InstanceID = viper.GetString("engine.instance_id")
	c.Engine.BufferSize = viper.GetInt("engine.buffer_size")
	c.Engine.Message = viper.GetString("engine.message")
}
