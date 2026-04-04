package config

import (
	"github.com/spf13/viper"
)

// FrequencyConfig holds configuration for a single frequency
type FrequencyConfig struct {
	BufferSize      int
	BatchSize       int
	FlushIntervalMs int
	Dispatcher      struct {
		MaxRetries  int
		BaseBackoff int
		MaxBackoff  int
	}
}

type Config struct {
	Instance struct {
		ID              string
		ProducerVersion string
	}
	Pipelines struct {
		EnabledFrequencies []string
	}
	FrequencyConfig map[string]*FrequencyConfig
	Users           struct {
		Count int
		Seed  int64
	}
	Signals struct {
		Enabled   bool
		Noise     map[string]any
		Anomalies struct {
			Enabled  bool
			Interval int
			Types    []string
		}
	}
	Chunking struct {
		Enabled        bool
		ChunkSizeBytes int
		Frequencies    []string
	}
	DLQ struct {
		Enabled   bool
		Type      string // "local" or "s3"
		Directory string
	}
	Kinesis struct {
		Enabled    bool
		StreamName string
		Region     string
	}
}

func (c *Config) Load() {
	viper.SetConfigName("config")
	viper.AddConfigPath("./config")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	// Load instance config
	c.Instance.ID = viper.GetString("instance.id")
	c.Instance.ProducerVersion = viper.GetString("instance.producer_version")

	// Load pipelines config
	c.Pipelines.EnabledFrequencies = viper.GetStringSlice("pipelines.enabled_frequencies")

	// Load per-frequency config
	c.FrequencyConfig = make(map[string]*FrequencyConfig)
	frequencies := []string{"second", "minute", "hour", "day"}
	for _, freq := range frequencies {
		freqCfg := &FrequencyConfig{
			BufferSize:      viper.GetInt("frequency_config." + freq + ".buffer_size"),
			BatchSize:       viper.GetInt("frequency_config." + freq + ".batch_size"),
			FlushIntervalMs: viper.GetInt("frequency_config." + freq + ".flush_interval_ms"),
		}
		freqCfg.Dispatcher.MaxRetries = viper.GetInt("frequency_config." + freq + ".dispatcher.max_retries")
		freqCfg.Dispatcher.BaseBackoff = viper.GetInt("frequency_config." + freq + ".dispatcher.base_backoff")
		freqCfg.Dispatcher.MaxBackoff = viper.GetInt("frequency_config." + freq + ".dispatcher.max_backoff")
		c.FrequencyConfig[freq] = freqCfg
	}

	// Load users config
	c.Users.Count = viper.GetInt("users.count")
	c.Users.Seed = int64(viper.GetInt("users.seed"))

	// Load signals config
	c.Signals.Enabled = viper.GetBool("signals.enabled")
	c.Signals.Noise = viper.GetStringMap("signals.noise")
	c.Signals.Anomalies.Enabled = viper.GetBool("signals.anomalies.enabled")
	c.Signals.Anomalies.Interval = viper.GetInt("signals.anomalies.interval")
	c.Signals.Anomalies.Types = viper.GetStringSlice("signals.anomalies.types")

	// Load chunking config
	c.Chunking.Enabled = viper.GetBool("chunking.enabled")
	c.Chunking.ChunkSizeBytes = viper.GetInt("chunking.chunk_size_bytes")
	c.Chunking.Frequencies = viper.GetStringSlice("chunking.frequencies")

	// Load DLQ config
	c.DLQ.Enabled = viper.GetBool("dlq.enabled")
	c.DLQ.Type = viper.GetString("dlq.type")
	c.DLQ.Directory = viper.GetString("dlq.directory")

	// Load Kinesis config
	c.Kinesis.Enabled = viper.GetBool("kinesis.enabled")
	c.Kinesis.StreamName = viper.GetString("kinesis.stream_name")
	c.Kinesis.Region = viper.GetString("kinesis.region")
}
