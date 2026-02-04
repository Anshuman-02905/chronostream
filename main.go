package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Hello World!!")

	for i := 0; i < 5; i++ {
		time.Sleep(1 * time.Second)
		current_time := time.Now()
		seconds := current_time.Second()
		minutes := current_time.Minute()
		hours := current_time.Hour()
		fmt.Println(seconds, minutes, hours)
	}
}
