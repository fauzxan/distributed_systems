package client

import (
	"core/message"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Client struct {
	Counter  int
	Lock     sync.Mutex
	Clientid int
	Channel  chan message.Message
}

// Call this by invoking client.send(msg, clientlist)
func (client *Client) Send(msg message.Message, clientlist *Clientlist, server *Server) {
	/*
		This function will simply send messages to the server, by first carrying out a couple of checks.
		The from id in the message should correspond to a valid client in the clientlist
		The client will acquire a lock to its own CS to update the counter value before sending, then append it to the message
		and send it to the server in one of its 5 channels.
		The channel to be sent in is chosen randomly.
	*/

	client.Lock.Lock()
	// clientlist.PrintClients()
	if clientlist.Check(msg.From) {
		client.Counter++
		fmt.Println("[", time.Now().UTC().String()[11:27], "] [Client Event] Sent from", client.Clientid, "to server. Sender counter is now", client.Counter) // Print the counter of the client id
		server.Channel[rand.Intn(5)] <- message.Message{
			Counter: client.Counter,
			Message: msg.Message,
			From:    client.Clientid,
		}

	} else {
		fmt.Println("[", time.Now().UTC().String()[11:27], "] [Client Event] [Error] Invalid send id", msg.From)
	}
	client.Lock.Unlock()
}

func (client *Client) Receive(clientlist *Clientlist) {
	/*
		This function will try and receive messages on behalf of the client. This will start an infinite loop for each client, listening
		to messages coming in from that channel. Upon reception,
		1. The client will increment its counter by comparing the counter in the message, and its own counter, then adding one to the max.
		2. It will print out a statement
	*/
	for i := 0; ; i++ {
		select {
		case msg := <-client.Channel:
			client.Lock.Lock()
			client.Counter = max(msg.Counter, client.Counter) + 1
			fmt.Println("[", time.Now().UTC().String()[11:27], "] [Client Event]", client.Clientid, "received a message from", msg.From, "Update counter value to", client.Counter)
			client.Lock.Unlock()
		default:
			time.Sleep(1 * time.Second)
			// fmt.Println("Client", client.Clientid , "is awake and waiting")

		}
	}
}
