package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Anshuman-02905/chronostream/internal/buffer"
	"github.com/Anshuman-02905/chronostream/internal/config"
	"github.com/Anshuman-02905/chronostream/internal/dispatcher"
	"github.com/Anshuman-02905/chronostream/internal/dlq"
	"github.com/Anshuman-02905/chronostream/internal/engine"
	"github.com/Anshuman-02905/chronostream/internal/event"
	"github.com/Anshuman-02905/chronostream/internal/monotime"
	"github.com/Anshuman-02905/chronostream/internal/scheduler"
	"github.com/Anshuman-02905/chronostream/internal/sequence"
	"github.com/Anshuman-02905/chronostream/internal/transport"
)

type ProducerStruct struct {
	version int
	uuid    uint8
}

type Event struct {
	uuid            int
	timestamp       time.Time
	FrequencyType   string
	monotonicSeqNum int

	Producer ProducerStruct
	Version  string
}

func main() {
	fmt.Println("Hello World!!")
	//First load the configuration
	var cfg config.Config
	cfg.Load()

	ts := monotime.RealTimeSource{}

	//buffer
	buf := buffer.New(10)

	//Sequencer
	seq := sequence.New()
	//scheduler
	sch := scheduler.New(event.FrequencySecond, &ts, cfg.Engine.BufferSize)

	dq, _ := dlq.NewFileDlq(cfg.DLQ.Directory, cfg.Engine.InstanceID, &ts)

	eng := engine.New(sch, seq, buf, cfg.Engine.ProducerVersion, cfg.Engine.InstanceID)
	ctx := context.Background()
	message := "Hello"
	go eng.Start(ctx, message)

	trans, err := transport.NewAwsKinesisTransport(ctx, cfg)
	if err != nil {
		panic(err)
	}
	defer trans.Close(ctx)

	disp := dispatcher.New(buf, trans, cfg, &ts, dq)

	go disp.StartBatch(ctx)
	select {}
}
