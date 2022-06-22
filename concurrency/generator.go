package concurrency

import (
	"fmt"
	"time"
)

func gen(msg string) <-chan string { // Returns recieve-only channel of strings
	c := make(chan string)

	go func() { // Launch gorutine from inside the function
		for i := 0; ; i++ {
			c <- fmt.Sprintf("%s %d", msg, i)
			time.Sleep(time.Duration(time.Second))
		}
	}()

	return c // Return channel to the caller
}

func Exec_gen() {
	c1 := gen("Generator 1")
	c2 := gen("Generator 2")

	for i := 0; i < 5; i++ {
		fmt.Printf("%q\n", <-c1)
		fmt.Printf("%q\n", <-c2)
	}

	fmt.Println("Finish")
}
