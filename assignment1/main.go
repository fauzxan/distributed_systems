package main

import (
	"core/client"
	"core/message"
	"fmt"
	"sync"
	"time"
)

func send(from int, clientlist *client.Clientlist, server *client.Server){
	go clientlist.Clients[from].Send(message.Message{Message: "Hello", From: from}, clientlist, server)
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

func createServer(clientlist *client.Clientlist) *client.Server{
	
	var chans [5]chan message.Message
	for i := range chans {
		chans[i] = make(chan message.Message)
	}
	return &client.Server{
		Counter: 0,
		Lock: sync.Mutex{},
		Serverid: 10000,
		Channel: chans,
		List: clientlist,
	}
}

func main(){
	// fmt.Println("yo")
	// Modify the line below to indicate the number of clients. 
	// START
	clientlist := createclients(10)
	// END
	server := createServer(clientlist) // creates 5 channels for the server and returns the address of the server. 
	
	
	fmt.Println("The server is active with id", server.Serverid)
	clientlist.PrintClients()

	// Modify the lines below to indicate which clients send to which clients
	// START
	send(0, clientlist, server)
	send(1, clientlist, server)
	send(2, clientlist, server)
	// END

	for _, client := range clientlist.Clients{
		go client.Receive(clientlist)
	}
	server.Receive()
	time.Sleep(5 * time.Second)
	clientlist.PrintClients()
}