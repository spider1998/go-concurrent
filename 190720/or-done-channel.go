package main

import (
	"fmt"
	"time"
)

func main() {

	ordone := func(done, c <-chan interface{}) <-chan interface{} {
		valStream := make(chan interface{})
		go func() {
			defer close(valStream)
			for {
				select {
				case <-done:
					return
				case v, ok := <-c:
					if !ok {
						return
					}
					select {
					case valStream <- v:
					case <-done:
					}
				}
			}
		}()
		return valStream
	}

	dosome := func(done, c chan interface{}) <-chan interface{} {
		someStream := make(chan interface{})
		go func() {
			defer close(someStream)
			for {
				select {
				case <-done:
					return
				case someStream <- 1:
				}
				time.Sleep(3 * time.Second)
			}
		}()
		return someStream
	}

	done := make(chan interface{})
	some := make(chan interface{})
	mychan := dosome(done, some)
	n := 1
	for val := range ordone(done, mychan) {
		fmt.Println(val.(int) + 1)
		n += 1
		if n == 5 {
			close(done)
		}
	}
	fmt.Println("all goroutine over!")

}
