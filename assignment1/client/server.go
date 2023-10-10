package client

import (
	"core/message"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Server struct {
	Counter  int
	Lock     sync.Mutex
	Serverid int
	Channel  [5]chan message.Message
	List     *Clientlist
}

type Queue struct {
	Lock  sync.Mutex
	Queue []message.Message
}

var queue Queue = Queue{
	sync.Mutex{},
	make([]message.Message, 0),
}

func (server *Server) Send(msg message.Message) {
	/*
		In this function, we will make use of a queue to send messages out to all the clients. Whenever some function calls server.Send(), we
		will add the message that is in the parameter to a queue, using a lock for the queue (because we don't want to modify the contents of
		the queue while it is being updated). This function will in turn call another function called queuesender(), which will iterate through
		all the elements in the queue structure, and send them in the FIFO order.
	*/
	queue.Lock.Lock()
	queue.Queue = append(queue.Queue, msg)
	fmt.Println("Message from", msg.From, "has been sent to the queue")
	queue.Lock.Unlock()
	server.queuesender()
}

func (server *Server) Receive() {
	/*
		This function will listen to messages on all the incoming channels of the server.
		In order to do this, it will infinitely listen to all the channels. Since we know, via implementation that the server only has 5 channels, we can just
		hardcode the 5 channels here to listen to them.
		Once we receive a message from one of the channels, server.Send() will be invoked after:
			1. updating the counter value of the server.
			2. appending the counter value to the message.
	*/
	for i := 0; ; i++ {
		/// Every message that you receive here must:
		// 1. Increment the counter value by comparing and retrieving the max. This comparison is to be done for receiving the messages. The counter needs to
		//		be further incremented before sending it.
		// 2. The calling the server.Send() function to send the values over to the queue.
		select {
		case msg := <-server.Channel[0]:
			server.Lock.Lock()
			server.Counter = max(server.Counter, msg.Counter) + 1
			message := message.Message{
				From:    msg.From,
				Counter: server.Counter,
				Message: msg.Message,
			}
			fmt.Println("Server channel", 0, "just received a message from", msg.From, ". Counter value is now", server.Counter)
			server.Lock.Unlock()
			server.Send(message)
		case msg := <-server.Channel[1]:
			server.Lock.Lock()
			server.Counter = max(server.Counter, msg.Counter) + 1
			message := message.Message{
				From:    msg.From,
				Counter: server.Counter,
				Message: msg.Message,
			}
			fmt.Println("Server channel", 1, "just received a message from", msg.From, ". Counter value is now", server.Counter)
			server.Lock.Unlock()
			server.Send(message)
		case msg := <-server.Channel[2]:
			server.Lock.Lock()
			server.Counter = max(server.Counter, msg.Counter) + 1
			message := message.Message{
				From:    msg.From,
				Counter: server.Counter,
				Message: msg.Message,
			}
			fmt.Println("Server channel", 2, "just received a message from", msg.From, ". Counter value is now", server.Counter)
			server.Lock.Unlock()
			server.Send(message)
		case msg := <-server.Channel[3]:
			server.Lock.Lock()
			server.Counter = max(server.Counter, msg.Counter) + 1
			message := message.Message{
				From:    msg.From,
				Counter: server.Counter,
				Message: msg.Message,
			}
			fmt.Println("Server channel", 3, "just received a message from", msg.From, ". Counter value is now", server.Counter)
			server.Lock.Unlock()
			server.Send(message)
		case msg := <-server.Channel[4]:
			server.Lock.Lock()
			server.Counter = max(server.Counter, msg.Counter) + 1
			message := message.Message{
				From:    msg.From,
				Counter: server.Counter,
				Message: msg.Message,
			}
			fmt.Println("Server channel", 4, "just received a message from", msg.From, ". Counter value is now", server.Counter)
			server.Lock.Unlock()
			server.Send(message)
		default:
			time.Sleep(1 * time.Second)
			fmt.Println("Server is alive and is listening to 5 channels")
		}
	}
}

func (server *Server) queuesender() {
	/*
		This function will iterate through all the elements of the queue, and send the messages in FIFO order to all the clients.
		The function performs the following:
			For each message in the queue-
				1. We will randomly calculate a 50% chance of sending the message. In the situation that we choose to send it, we will follow steps 2-4. Else 5.
				2. For every message in the message queue, we will compare the counter in the message to server.Counter. The maximum between the two will be appended to
					the message.
				3. The message will now be sent to the all the channels in the clientlist, except for the one with clientid == msg.From.
	*/

	for _, msg := range queue.Queue {
		// queue.Lock.Lock()
		if rand.Intn(10) < 5 {
			server.Lock.Lock()
			for id, client := range server.List.Clients {
				server.Counter = max(server.Counter, msg.Counter) + 1
				// client.Lock.Lock()
				if id == msg.From {
					continue
				}
				client.Channel <- message.Message{
					From:    msg.From,
					Counter: server.Counter,
					Message: msg.Message,
				}
				fmt.Println("[", time.Now().UTC(), "] [Server Event] Server has forwarded the message from", msg.From, "to", id)
				// client.Lock.Unlock()
			}
			server.Lock.Unlock()
		} else {
			fmt.Println("[", time.Now().UTC(), "] [Server Event] Server has dropped the message from", msg.From)
		}
		queue.Lock.Lock()
		if len(queue.Queue) > 1 {
			queue.Queue = queue.Queue[1:]
		} else {
			queue.Queue = make([]message.Message, 0)
		}
		queue.Lock.Unlock()
	}
}

// func (server *Server) sendtoclient(client *Client){
// 	client.Lock.Lock()
// 	defer client.Lock.Unlock()
// 	client.Channel <-
// }
