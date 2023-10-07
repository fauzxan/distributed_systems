package main

import (
	"core/client"
	"core/message"
	"fmt"
	"sync"
	"time"
)

func send(from int, to int, clientlist *client.Clientlist){
	go clientlist.Clients[from].Send(message.Message{Message: "Hello", To: to}, clientlist)
}


func createclients(num int) *client.Clientlist{
	clientlist := client.Clientlist{
		Lock: sync.Mutex{},
		Clients: make(map[int]*client.Client),
	}
	for i:=0;i<num;i++{
		clientlist.Add( &client.Client{
			Counter: 0,
			Lock: sync.Mutex{},
			Clientid: i,
			Channel: make(chan message.Message),
		})
	}
	return &clientlist
}

func main(){
	// fmt.Println("yo")
	// Modify the line below to indicate the number of clients. 
	// START
	clientlist := createclients(10)
	// END
	
	fmt.Println("The server is", clientlist.Getmax())
	clientlist.PrintClients()

	// Modify the lines below to indicate which clients send to which clients
	// START
	send(0, 1, clientlist)
	send(1, 0, clientlist)
	send(2, 9, clientlist)
	// END

	for _, client := range clientlist.Clients{
		go client.Receive(clientlist)
	}
	time.Sleep(5 * time.Second)
	clientlist.PrintClients()
}