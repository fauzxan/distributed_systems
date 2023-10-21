package client

import (
	"core/message"
	"math/rand"
	"sync"
	"time"

	"github.com/fatih/color"
)

type Server struct {
	Counter  []int
	Lock     sync.Mutex
	Serverid int // is equal to the highest id in the clientlist. By design of this implementation.
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

var svr = color.New(color.FgWhite).Add(color.BgBlack)
var q = color.New(color.FgCyan).Add(color.BgBlack)
var causality = color.New(color.FgHiRed).Add(color.BgBlack)

func (server *Server) Send(msg message.Message) {
	/*
		In this function, we will make use of a queue to send messages out to all the clients. Whenever some function calls server.Send(), we
		will add the message that is in the parameter to a queue, using a lock for the queue (because we don't want to modify the contents of
		the queue while it is being updated). This function will in turn call another function called queuesender(), which will iterate through
		all the elements in the queue structure, and send them in the FIFO order.
	*/
	queue.Lock.Lock()
	queue.Queue = append(queue.Queue, msg)
	svr.Println("[", time.Now().UTC().String()[11:27], "] [Server Event] Message from", msg.From, "has been sent to the queue")
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
			server.Counter = server.UpdateVectorClock(msg.Counter, server.Counter, msg.From)
			server.Counter[server.Serverid]++
			message := message.Message{
				From:    msg.From,
				Counter: server.Counter,
				Message: msg.Message,
			}
			svr.Println("[", time.Now().UTC().String()[11:27], "] [Server Event] Server channel", 0, "just received a message from", msg.From, "| Counter value is now", server.Counter)
			go server.Send(message)
		case msg := <-server.Channel[1]:
			server.Counter = server.UpdateVectorClock(msg.Counter, server.Counter, msg.From)
			server.Counter[server.Serverid]++
			message := message.Message{
				From:    msg.From,
				Counter: server.Counter,
				Message: msg.Message,
			}
			svr.Println("[", time.Now().UTC().String()[11:27], "] [Server Event] Server channel", 1, "just received a message from", msg.From, "| Counter value is now", server.Counter)
			go server.Send(message)
		case msg := <-server.Channel[2]:
			server.Counter = server.UpdateVectorClock(msg.Counter, server.Counter, msg.From)
			server.Counter[server.Serverid]++
			message := message.Message{
				From:    msg.From,
				Counter: server.Counter,
				Message: msg.Message,
			}
			svr.Println("[", time.Now().UTC().String()[11:27], "] [Server Event] Server channel", 2, "just received a message from", msg.From, "| Counter value is now", server.Counter)
			go server.Send(message)
		case msg := <-server.Channel[3]:
			server.Counter = server.UpdateVectorClock(msg.Counter, server.Counter, msg.From)
			server.Counter[server.Serverid]++
			message := message.Message{
				From:    msg.From,
				Counter: server.Counter,
				Message: msg.Message,
			}
			svr.Println("[", time.Now().UTC().String()[11:27], "] [Server Event] Server channel", 3, "just received a message from", msg.From, "| Counter value is now", server.Counter)
			go server.Send(message)
		case msg := <-server.Channel[4]:
			server.Counter = server.UpdateVectorClock(msg.Counter, server.Counter, msg.From)
			server.Counter[server.Serverid]++
			message := message.Message{
				From:    msg.From,
				Counter: server.Counter,
				Message: msg.Message,
			}
			svr.Println("[", time.Now().UTC().String()[11:27], "] [Server Event] Server channel", 4, "just received a message from", msg.From, "| Counter value is now", server.Counter)
			go server.Send(message)
		default:
			time.Sleep(1 * time.Second)
			svr.Println("[", time.Now().UTC().String()[11:27], "] [Server Ping] Server is alive and is listening to 5 channels")
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
	queue.Lock.Lock()
	defer queue.Lock.Unlock()
	if len(queue.Queue) == 0 {
		queue.Queue = make([]message.Message, 0)
		return
	}
	for _, msg := range queue.Queue {
		q.Println("[Queue state] Queue is sending now... Head of queue is from", queue.getQueueIds())
		if rand.Intn(10) < 5 {
			for id, client := range server.List.Clients {
				server.Counter[server.Serverid]++
				if id == msg.From {
					continue
				}
				client.Channel <- message.Message{
					From:    msg.From,
					Counter: server.Counter,
					Message: msg.Message,
				}
				svr.Println("[", time.Now().UTC().String()[11:27], "] [Server Event] Server has forwarded the message from", msg.From, "to", id, "| Counter value is now", server.Counter)
			}
		} else {
			svr.Println("[", time.Now().UTC().String()[11:27], "] [Server Event] Server has dropped the message from", msg.From)
		}
		if len(queue.Queue) > 1 {
			queue.Queue = queue.Queue[1:]
		} else {
			queue.Queue = make([]message.Message, 0)
		}
		q.Println("[Queue state] Removing from queue now... State of queue is", queue.getQueueIds())
	}
}

/*
*****************
UTILITY FUNCTIONS
*****************
*/
func (server *Server) CloseChannels() {
	for idx, channel := range server.Channel {
		close(channel)
		svr.Println("[", time.Now().UTC().String()[11:27], "] [Server Event] Server channel", idx, "closed")
	}
}

func (server *Server) UpdateVectorClock(V1 []int, V2 []int, from int) []int {
	/*
		Does the same thing as the one in vectorclock.go
		V1 is the clock in the message
		V2 is the server clock.
	*/
	if (V2[from] > V1[from]){
		causality.Println("Potential causality violation")
	}
	for i := range V1 {
		// if (V2[i] > V1[i] && i == from) {
		// 	// there is a potential causality violation
		// 	causality.Println("There is a potential causality violation!!", "Incoming value at position", i, ":", V1[i], "Local value:", V2[i])
		// }
		V2[i] = max(V2[i], V1[i])
	}
	return V2
}

func (queue *Queue) getHead() int { // that is infact what she said
	if len(queue.Queue) != 0 {
		return queue.Queue[0].From
	}
	return -1
}

func (queue *Queue) getQueueIds() []int {
	var list = make([]int, len(queue.Queue))
	for i, msg := range queue.Queue {
		list[i] = msg.From
	}
	return list
}
