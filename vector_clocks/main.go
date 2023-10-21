package main

import (
	"core/client"
	"core/message"
	"fmt"
	// "os"
	// "os/signal"
	"sync"
	// "syscall"
	"time"
)

func send(from int, clientlist *client.Clientlist, server *client.Server) {
	message := message.Message{Message: "Hello", From: from}
	go clientlist.Clients[from].Send(message, clientlist, server)
}

func createclients(num int) *client.Clientlist {
	clientlist := client.Clientlist{
		Lock:    sync.Mutex{},
		Clients: make(map[int]*client.Client),
	}
	
	for i := 0; i < num; i++ {
		var vector = make([]int, num+1)
		for i:=0;i<len(vector);i++{
			vector[i] = 0
		}
		clientlist.Add(&client.Client{
			Id: i,
			Counter:  vector,
			Lock:     sync.Mutex{},
			Clientid: i,
			Channel:  make(chan message.Message),
		})
	}
	return &clientlist
}

func createServer(clientlist *client.Clientlist, num int) *client.Server {

	var chans [5]chan message.Message
	for i := range chans {
		chans[i] = make(chan message.Message)
	}
	var vector = make([]int, num+1)
	for i:=0;i<len(vector);i++{
		vector[i] = 0
	}
	return &client.Server{
		Counter: vector,
		Lock:     sync.Mutex{},
		Serverid: num,
		Channel:  chans,
		List:     clientlist,
	}
}


func main() {

	// Modify the line below to indicate the number of clients.
	// START
	num := 10 // number of clients
	// END	
	clientlist := createclients(num)
	server := createServer(clientlist, num) // creates 5 channels for the server and returns the address of the server.
	SetupCloseHandler(server, clientlist)

	fmt.Println("[", time.Now().UTC().String()[11:27], "] [Server Event] The server is active with id", server.Serverid)
	clientlist.PrintClients()

	// Modify the lines below to indicate which clients send to which clients
	// START
	send(0, clientlist, server)
	send(0, clientlist, server)
	send(0, clientlist, server)
	send(0, clientlist, server)
	send(0, clientlist, server)
	send(1, clientlist, server)
	send(2, clientlist, server)
	send(3, clientlist, server)
	send(4, clientlist, server)
	send(5, clientlist, server)
	// END

	for _, client := range clientlist.Clients {
		go client.Receive(clientlist)
	}
	go server.Receive()

	// Keep the parent thread alive!!
	for {
		time.Sleep(1000)
	}	
}
