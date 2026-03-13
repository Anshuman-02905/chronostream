package config

import (
	"strconv"

	"github.com/spf13/viper"
)

type Config struct {
	Engine struct {
		ProducerVersion string
		InstanceID      string
		BufferSize      int
		Message         string
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
	val, _ := strconv.Atoi(viper.GetString("engine.buffer_size"))
	c.Engine.BufferSize = val
	c.Engine.Message = viper.GetString("engine.message")
}
