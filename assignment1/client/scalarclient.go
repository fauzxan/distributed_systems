package client

import (
	"fmt"
	"sync"
	"core/message"
	"math/rand"
)

type Client struct{
	Counter int
	Lock sync.Mutex
	Clientid int
	Channel chan message.Message
}


// Call this by invoking client.send(msg, clientlist)
func (client *Client) Send(msg message.Message, clientlist *Clientlist) {
	client.Lock.Lock()
	serverid := clientlist.Getmax()
	fmt.Println("Inside the sender", serverid)
	if (clientlist.Check(msg.To) && client.Clientid != serverid){ // If the sender is not the server, send it to the server
		client.Counter++ // Increment the event counter of self
		clientlist.Clients[serverid].Channel <- message.Message{
			Counter: client.Counter, //  always needs to be the counter of the sender. 
			Message: msg.Message,
			From: client.Clientid, // Keep the from and to static because the receiver should ultimately know who actually sent the message.
			To: msg.To,
		} // Get the server channel, and send the message
		fmt.Println("Sent from", client.Clientid, "to server. Sender counter is now", client.Counter) // Print the counter of the client id
	}else if(clientlist.Check(msg.To) && client.Clientid == serverid){// If the sender is the server, then send it to the client it is intended for, and also replace the counter in the message to the server counter
		if (rand.Intn(10) < 5){
			client.Counter++
			clientlist.Clients[msg.To].Channel <- message.Message{
				Counter: client.Counter, // always needs to tbe the coutner of the sender. 
				Message: msg.Message,
				From: msg.From, // Keep the from and to static because the receiver should ultimately know who actually sent the message.
				To: msg.To,
			}
		}
	}
	client.Lock.Unlock()
}

func receive() {}
