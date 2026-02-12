package main

import (
	"fmt"
	"time"

	"github.com/Anshuman-02905/chronostream/internal/monotime"
)

type ProducerStruct struct {
	version int
	uuid    int
}

type Event struct {
	uuid            int
	timestamp       time.Time
	FrequencyType   string
	monotonicSeqNum int
	//Deterministic Signal Seed
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
