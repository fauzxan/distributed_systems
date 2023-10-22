# Lamport's scalar clock [1.1 and 1.2]

## Introduction
The files under ~/lamport_clock combines the deliverables for 1.1 and 1.2:
1. Setup up a client and a server that is able to meet the requirements of the system described in question 1
2. Introduce some logic that is able to bring about total ordering of all the events. 

In order to run this file, go to the parent directory, run `go build`. This will generate a file called `core`. You can run this file as an executable - `./core`.

### Editable sections
The file name `main.go` has the main function with some sections clearly demarcated to indicate that you can edit them. 
```go
// START
<modifiable section>
//END
```
To be able to observe the output, some starter send messages have already been included. 
Periodic pinging from all the nodes can be toggled on by un-commenting the following code in `func main()`:
```go
enablePeriodicPinging() // <-- uncomment this line to enable periodic pinging from all the clients.
```

# Design:

The main client logic for the code has been written in `./lamport_clock/client/client.go` file. It does the following:
1. `client.Send(msg, clientlist, server)` : When called, will check the `clientlist` and send the `msg` to the `server` on one of its channels. The logic for choosing the channel is to choose a random channel.
2. `client.Receive(clientlist)` : This runs an infinite loop listening to messages on its channel. Each client in the clientlist has its own channel. When sending a message to a specific client, the server will use the channel of the client to send it a message. Upon receiving, the client will update it's counter based on lamport's logic -> max(incoming counter, local counter) + 1. Then it will put the sender id in a queue. This queue will be printed out in the end to indicate the order of messages received by all the clients.

The main server logic for the code has been written in `./lamport_clock/client/server.go` file. It does the following:
1. `server.Receive()`: Runs an infinite loop listening to 5 channels. Upon receiving a message, it updates counter using lamport's logic, updates the counter value in the message, and calls server.Send()
2. `server.Send()`: Puts the message in a queue, and calls server.queuesender()
3. `server.queuesender()`: will run through the queue and send the messages in FIFO order to all the clients. This process also locks the queue when sending messages, so that the queue doesn't get modified from different threads during the process of sending. The queue is also designed to randomly drop messages with a 50% chance. This is to simulate send failure. 

Main logic on the server side:
```go
for _, msg := range queue.Queue {
		if rand.Intn(10) < 5 {
            // code to extract the message from the queue, and send it into the channel of the client
            ...
        }
}
```

From, the client side, the logic is simple: 

Whenever, client.send() is called, it will forward the message into the server's channel. The sever logic will be executed as above^

# Testing and Observations:
Running the code using `.\core` will generate logs that you can use to verify the total ordering of the messages. 

**The messages have been colored to enhance readability of the logs.**

> Green messages signify `client` events

> Blue messages signify queue states before and after `queuesender()` is called

> White messages signify server events

> Occassionally, some messages are uncolored due to bugs in the color package imported 😭😭 Please have a look at these messages to see if you missed something important.

The main.go has the following send commands, for this test:
```go
	send(0, clientlist, server)
	send(1, clientlist, server)
	send(2, clientlist, server)
	send(3, clientlist, server)
	send(4, clientlist, server)
	send(5, clientlist, server)
```
The order in which they are sent is not enforced, as they are all go routines, but the order in which the server receives them and sends them out is enforced.
Upon executing `./core`, you will see the following:
![image](https://github.com/fauzxan/distributed_systems/assets/92146562/497c4331-331c-401f-97af-9f0847da658e)
As message from 1 was received by the server first, it is being sent into the queue. We can see the state of the queue in blue as [1]. The server then starts sending the message from 1 to all the clients. We can see the logs printed out by the client as follows:
![image](https://github.com/fauzxan/distributed_systems/assets/92146562/7bccf70f-feb7-43c3-b1ae-435e02ad97aa)
From the above screenshot, we can also see that the server received a message from 5.

After sending to all the clients, the queue pops its head, and logs the state of the queue, which is now [], or empty.
Then the message received second, 0 is added to the queue, and is sent out:
![image](https://github.com/fauzxan/distributed_systems/assets/92146562/c3369613-2af0-4e78-9fee-421064a79ac8)

After this, the server seems to have received messages from 2,4,3 and 5 at the same time. So the state of the queue is [2,4,3,5]. It then starts sending the messages from 2 and 4:
![image](https://github.com/fauzxan/distributed_systems/assets/92146562/7c271245-2ac1-4241-a58b-5b004a27b502)

After sending 2 and 4, the message from 3 gets dropped and the message from 5 doesn't get dropped, and is sent, as shown:
![image](https://github.com/fauzxan/distributed_systems/assets/92146562/49939783-7ffa-4772-ad3d-08c0378cab68)

Once you see the following output on the terminal, it means that the clients have stopped sending all together:
![image](https://github.com/fauzxan/distributed_systems/assets/92146562/92ab8dd0-2dd0-4acc-9498-f380275ec75f)

You can now press ctrl+C to see the order of messages received by the different clients:
![image](https://github.com/fauzxan/distributed_systems/assets/92146562/06b7299f-d072-4976-bbc0-0dcb88e63119)

From the screenshot above, we can see that the messages are received in the same order in all the clients. 
We can also see that the clients don't receive messages from themselves. As per the screenshot, the universal order of messages received is: [1, 0, 2, 4, 5]

# Answers 
### 1.1
Clients are able to send messages to the server, which then flips a coin to see if the message should be forwarded or not. 
### 1.2
There is a total ordering of messages sent out from the server, as the server maintains a queue and orders the messages to be sent out. So all the clients receive messages in the same order. This doesn't fully mitigate the causality violation possibility all together, as the order of receiving messages at the server is still not the same as the order in which the clients sent out the messages. 
