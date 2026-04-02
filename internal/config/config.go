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

		BatchSize     int
		FlushInterval int //ms
	}
	DLQ struct {
		Enabled   bool
		Directory string
	}
	Kinesis struct {
		Enabled      bool
		StreamName   string
		Region       string
		MaxRecordAge int //ms before retry
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

	c.Dispatcher.MaxRetries = viper.GetInt("dispatcher.max_retries")
	c.Dispatcher.BaseBackoff = viper.GetInt("dispatcher.base_backoff")
	c.Dispatcher.MaxBackoff = viper.GetInt("dispatcher.max_backoff")
	c.Dispatcher.BatchSize = viper.GetInt("dispatcher.batch_size")
	c.Dispatcher.FlushInterval = viper.GetInt("dispatcher.flush_interval")

	c.DLQ.Enabled = viper.GetBool("dlq.enabled")
	c.DLQ.Directory = viper.GetString("dlq.directory")

	c.Kinesis.Enabled = viper.GetBool("kinesis.enabled")
	c.Kinesis.StreamName = viper.GetString("kinesis.stream_name")
	c.Kinesis.Region = viper.GetString("kinesis.region")
	c.Kinesis.MaxRecordAge = viper.GetInt("kinesis.max_record_age")
}
