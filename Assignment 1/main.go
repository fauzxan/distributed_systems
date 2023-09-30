package main

import (
	"fmt"
	"time"
)



func sender(ch chan<- Message, from int, to int) {
	fmt.Println("Sending message from", from, "to", to)
	time.Sleep(200 * time.Millisecond) // Simulate some work for 200 milliseconds
	ch <- Message{id: 1,from: from, to: to, selfcounter: 0}
}

func receiver(ch <-chan Message) {

	select {
	case msg := <-ch:
		fmt.Println("Received:", msg)
	case <-time.After(300 * time.Millisecond):
		fmt.Println("Timeout! No message received.")
	}
}

func main() {
	channel := make(chan Message)
	
	go sender(channel, 1, 2) // Send a message after 200 milliseconds
	go sender(channel, 1, 2) // Send a message after 200 milliseconds
	go sender(channel, 1, 2) // Send a message after 200 milliseconds
	go sender(channel, 1, 2) // Send a message after 200 milliseconds

	go receiver(channel)               // Receive a message after 100 milliseconds

	time.Sleep(500 * time.Millisecond) // Wait for the goroutines to finish
}