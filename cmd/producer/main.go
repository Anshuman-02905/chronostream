package main

import (
	"fmt"
	"time"

	"github.com/Anshuman-02905/chronostream/internal/monotime"
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
	ticker := time.NewTicker(500 * time.Millisecond)
	ts := monotime.RealTimeSource{}
	for range ticker.C {
		fmt.Println(ts.Now())
	}
}
