package main

import (
	"fmt"
	"time"

	"github.com/Anshuman-02905/chronostream/internal/monotime"
)

func main() {
	fmt.Println("Hello World!!")
	ticker := time.NewTicker(500 * time.Millisecond)
	ts := monotime.RealTimeSource{}
	for range ticker.C {
		fmt.Println(ts.Now())
	}

}
