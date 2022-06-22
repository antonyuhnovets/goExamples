package concurrency

import (
	"fmt"
	"time"
)

func myFunc(done chan string) {
	for i := 0; i < 10; i++ {
		time.Sleep(time.Millisecond * 500)
		fmt.Println(i, " myFunc")
	}
	fmt.Println("finished loop in myFunc")
	done <- "goroutine finished" // send the message into the channel
}

func Exec() {
	done := make(chan string) // make the "done" channel
	go myFunc(done)           // run myFunc on a goroutine
	for i := 0; i < 5; i++ {
		time.Sleep(time.Millisecond * 500)
		fmt.Println(i, " main")
	}
	msg := <-done // receive from the channel
	fmt.Println(msg)
}
