package client

import (
	"core/message"
	"fmt"
	"math/rand"
	"sync"
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
	// clientlist.PrintClients()
	serverid := clientlist.Getmax()
	if (clientlist.Check(msg.To) && client.Clientid != serverid){ // If the sender is not the server, send it to the server
		client.Counter++ // Increment the event counter of self
		clientlist.Clients[serverid].Channel <- message.Message{ // Get the server channel, and send the message
			Counter: client.Counter, //  always needs to be the counter of the sender. 
			Message: msg.Message,
			From: client.Clientid, // Keep the from and to static because the receiver should ultimately know who actually sent the message.
			To: msg.To,
		} 
		fmt.Println("Sent from", client.Clientid, "to server. Sender counter is now", client.Counter) // Print the counter of the client id
	}else if(clientlist.Check(msg.To) && client.Clientid == serverid){// If the sender is the server, then send it to the client it is intended for, and also replace the counter in the message to the server counter
		if (rand.Intn(10) < 5){
			// client.Counter++ // Don't need this as the server messages are 
			clientlist.Clients[msg.To].Channel <- message.Message{
				Counter: client.Counter, // always needs to tbe the coutner of the sender. 
				Message: msg.Message,
				From: msg.From, // Keep the from and to static because the receiver should ultimately know who actually sent the message.
				To: msg.To,
			}
			fmt.Println("The server forwarded the message from", msg.From, "Counter value of server is", client.Counter)
		}else{
			fmt.Println("The server just dropped the message from", msg.From)
		}
	}
	client.Lock.Unlock()
}

func (client *Client) Receive(clientlist *Clientlist) {
	for i:=0;;i++{
		select {
		case msg := <- client.Channel:
			
			if client.Clientid == msg.To{ // If I am the intended receiver, I will receive the message and increment my counter. 
				client.Lock.Lock()
				var sendercounter int
				if (msg.To == clientlist.Getmax()){
					sendercounter = clientlist.Clients[msg.From].Counter
				}else{
					sendercounter = clientlist.Clients[clientlist.Getmax()].Counter
				}
				
				client.Counter = max(sendercounter, client.Counter) + 1
				fmt.Println("Client", client.Clientid, "received a message from", msg.From, "Update counter value to", client.Counter)
				client.Lock.Unlock()
			}else{// If I am not the intended receiver, I must be the server, hence I will forward the message to the intended msg.To
				client.Lock.Lock()
				sendercounter := clientlist.Clients[msg.From].Counter
				client.Counter = max(sendercounter, client.Counter) + 1
				client.Lock.Unlock()
				client.Send(message.Message{
					From: msg.From,
					To: msg.From,
					Counter: client.Counter,
					Message: msg.Message,
				}, clientlist)
				
			}	
		}
	}
}
