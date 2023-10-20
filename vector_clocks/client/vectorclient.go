package client

import (
	"core/message"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Client struct {
	Id int
	Counter  []int
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
		client.Counter[client.Id]++
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
			client.Counter = client.UpdateVectorClock(msg.Counter, client.Counter)
			client.Counter[client.Clientid] ++
			fmt.Println("[", time.Now().UTC().String()[11:27], "] [Client Event]", client.Clientid, "received a message from", msg.From, "Update counter value to", client.Counter)
			client.Lock.Unlock()
		default:
			time.Sleep(1 * time.Second)
			// fmt.Println("Client", client.Clientid , "is awake and waiting") // health check for each client process

		}
	}
}


/*
************************
UTILITY FUNCTION
************************
*/
func (client *Client) UpdateVectorClock(V1 []int, V2 []int) []int{
	/*
	Input:
		V1 is the incoming vector.
		V2 is the vector on the local process. 

	This function receives two vectors- one from the message that just came in, and one from the local process.
	It will compare all the elements in one vector to all the elements in the other vector, and:
	1. See if there are any causality violations
		If any element in V2 is greater than any element in V2, then there is a causality violation. In this case return []. The receiver will receive this and see that there
		is a causality violation, and flag it to the terminal. It will still continue running, but there is a potential causality violation. 
	2. Update each element in V1 and V2 as max(V1.1, V2.1), max(V1.2, V2.2) and so on

	Finally, this function should return the updated V2, which will be places into the receiver process
	*/
	for i := range V1{
		if (V2[i] > V1[i] && i != client.Clientid){
			// there is a potential causality violation
			// return make([]int, 0)
			fmt.Println("There is a potential causality violation!!")
		}
		V2[i] = max(V2[i], V1[i])
	}
	return V2
}