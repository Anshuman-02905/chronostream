package pipeline

import (
	"context"
	"fmt"

	"github.com/Anshuman-02905/chronostream/internal/config"
	"github.com/Anshuman-02905/chronostream/internal/dispatcher"
	"github.com/Anshuman-02905/chronostream/internal/event"
	"github.com/Anshuman-02905/chronostream/internal/monotime"
	"github.com/Anshuman-02905/chronostream/internal/transport"
	"github.com/Anshuman-02905/chronostream/internal/user"
)

type PipelineGroup struct {
	pipelines map[event.Frequency]*FrequencyPipeline
}

// NewGroup creates a PipelineGroup by dynamically iterating over
// cfg.Pipelines.EnabledFrequencies and building a FrequencyPipeline
// for each one using the per-frequency config from cfg.FrequencyConfig.
func NewGroup(cfg config.Config, tsp transport.Transport) (*PipelineGroup, error) {
	pipelines := make(map[event.Frequency]*FrequencyPipeline)
	ts := &monotime.RealTimeSource{}

	// Create a single registry shared across all frequency pipelines
	// All 4 frequencies simulate the same set of users (count and seed from config)
	registry, err := user.NewUserRegistry(cfg.Users.Count, cfg.Users.Seed)
	if err != nil {
		return nil, fmt.Errorf("failed to create user registry: %w", err)
	}

	for _, freqStr := range cfg.Pipelines.EnabledFrequencies {
		// Convert config string ("second") to typed enum (FrequencySecond)
		freq, err := event.ParseFrequency(freqStr)
		if err != nil {
			return nil, fmt.Errorf("invalid frequency in config: %w", err)
		}

		// Lookup per-frequency config
		freqCfg := cfg.FrequencyConfig[freqStr]
		if freqCfg == nil {
			return nil, fmt.Errorf("frequency_config not found for %q", freqStr)
		}

		// Build a PipelineConfig from the per-frequency settings
		pCfg := PipelineConfig{
			Frequency:       freq,
			BufferSize:      freqCfg.BufferSize,
			DLQDirectory:    cfg.DLQ.Directory,
			InstanceID:      cfg.Instance.ID,
			ProducerVersion: cfg.Instance.ProducerVersion,
			TimeSource:      ts,
			Users:           registry,
			Dispatcher: dispatcher.DispatcherConfig{
				MaxRetries:    freqCfg.Dispatcher.MaxRetries,
				BaseBackoff:   freqCfg.Dispatcher.BaseBackoff,
				MaxBackoff:    freqCfg.Dispatcher.MaxBackoff,
				BatchSize:     freqCfg.BatchSize,
				FlushInterval: freqCfg.FlushIntervalMs,
			},
		}

		p, err := New(pCfg, tsp)
		if err != nil {
			return nil, fmt.Errorf("failed to create %s pipeline: %w", freqStr, err)
		}
		pipelines[freq] = p
		fmt.Printf("Created %s pipeline (buffer=%d, retries=%d)\n",
			freqStr, freqCfg.BufferSize, freqCfg.Dispatcher.MaxRetries)
	}

	return &PipelineGroup{
		pipelines: pipelines,
	}, nil
}

// StartAll starts all frequency pipelines concurrently
// Each pipeline runs its engine and dispatcher in separate goroutines
func (pg *PipelineGroup) StartAll(ctx context.Context) {
	for freq, p := range pg.pipelines {
		fmt.Printf("Starting %v frequency pipeline\n", freq)
		p.Start(ctx, "Deterministic Pulse")
	}
}

// StopAll stops all frequency pipelines gracefully
func (pg *PipelineGroup) StopAll() {
	for freq, p := range pg.pipelines {
		fmt.Printf("Stopping %v frequency pipeline\n", freq)
		p.Stop()
	}
}

// Status returns the status of each frequency pipeline
func (pg *PipelineGroup) Status() map[event.Frequency]PipelineStatus {
	statuses := make(map[event.Frequency]PipelineStatus)
	for freq, p := range pg.pipelines {
		p.statusMutex.RLock()
		statuses[freq] = p.status
		p.statusMutex.RUnlock()
	}
	return statuses
}
