package main

import (
	"context"
	"fmt"

	"github.com/Anshuman-02905/chronostream/internal/config"
	"github.com/Anshuman-02905/chronostream/internal/pipeline"
	"github.com/Anshuman-02905/chronostream/internal/transport"
)

func main() {
	fmt.Println("Hello World!!")

	// Load configuration from config.yaml
	var cfg config.Config
	cfg.Load()
	fmt.Printf("Loaded config: instance_id=%s, enabled_frequencies=%v\n",
		cfg.Instance.ID, cfg.Pipelines.EnabledFrequencies)

	// Create Kinesis transport (shared across all frequency pipelines)
	ctx := context.Background()
	trans, err := transport.NewAwsKinesisTransport(ctx, cfg)
	if err != nil {
		panic(err)
	}
	defer trans.Close(ctx)

	// Create PipelineGroup with all enabled frequency pipelines
	// This replaces the single-frequency setup above
	group, err := pipeline.NewGroup(cfg, trans)
	if err != nil {
		panic(err)
	}

	// Start all frequency pipelines concurrently
	// Each frequency (Second, Minute, Hour, Day) runs independently
	group.StartAll(ctx)
	fmt.Println("All frequency pipelines started")

	// Keep producer running indefinitely
	select {}
}
