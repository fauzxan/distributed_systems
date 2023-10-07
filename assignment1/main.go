package main

import (
	"core/client"
	"core/message"
	"fmt"
	"sync"
	"time"
)

func send(from int, to int, clientlist *client.Clientlist){
	go clientlist.Clients[from].Send(message.Message{Message: "Hello", To: to}, &client.Clientlist{})
}

func main(){
	// fmt.Println("yo")
	clientlist := client.Clientlist{
		Lock: sync.Mutex{},
		Clients: make(map[int]*client.Client),
	}
	for i:=0;i<10;i++{
		clientlist.Add( &client.Client{
			Counter: 0,
			Lock: sync.Mutex{},
			Clientid: i,
			Channel: make(chan message.Message),
		})
	}
	fmt.Println("The server is", clientlist.Getmax())
	clientlist.PrintClients()

	send(0, 1, &clientlist)
	send(1, 0, &clientlist)
	time.Sleep(5 * time.Second)
	clientlist.PrintClients()
}