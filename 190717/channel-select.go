package main

import (
	"fmt"
	"time"
)

func main() {
	var list chan int
	list = make(chan int, 100)
	go func() {
		for i := range list {
			fmt.Print(i)
		}
	}()

	go func() {
		select {
		case <-list:
			fmt.Println("insert data ")
		default:
			fmt.Println("close....")
		}
	}()

	go func() {
		select {}
		fmt.Println("waite.....")
	}()

	for {
		list <- 1
		time.Sleep(time.Second * 5)
	}
}
