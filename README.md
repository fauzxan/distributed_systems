🚀Golang implementation of lamport's logical clock, vector clock, as well as bully algorithm.


# Part 1

Please read the following instructions to see how to run this file, and how to make amendments to different parts of the file to observe different scenarios. The explicit answers have been given below the `Testing and Observations` section. 

The files under ~/lamport_clock combines the deliverables for 1.1 and 1.2:
1. Setup up a client and a server that is able to meet the requirements of the system described in part 1
2. Introduce some logic that is able to bring about total ordering of all the events. 

To run this code, go to the root directory of this part - `~/lamport_clock` and run the following command:
```shell
go build && ./core
```

### Editable sections
The file name `main.go` has the main function with some sections clearly demarcated to indicate that you can edit them. 
```go
// START
<modifiable section>
//END
```
To be able to observe the output, some starter send messages have already been included. This is to enable determinism, and allow us to read the output in a slow, and determinisitc way. If you can observe log outputs at very high speeds, then please uncomment the enablePeriodicPinging() method as shown below. Periodic pinging is part of the requirements for this part. 
Periodic pinging from all the nodes can be toggled on by un-commenting the following code in `func main()`:
```go
go enablePeriodicPinging() // <-- uncomment this line to enable periodic pinging from all the clients.
```
The periodic pinging also proves that the system works with atleast 10 clients. You can increase the number of clients by changing 

### Design:

The main client logic for the code has been written in `./lamport_clock/client/client.go` file. It does the following:
1. `client.Send(msg, clientlist, server)` : When called, will check the `clientlist` and send the `msg` to the `server` on one of its channels. The logic for choosing the channel is to choose a _random_ channel.
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

### Testing and Observations:
Running the code using `.\core` will generate logs that you can use to verify the total ordering of the messages. 

**The messages have been colored to enhance readability of the logs.**

> Green messages signify `client` events

> Blue messages signify queue states before and after `queuesender()` is called

> White messages signify server events


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

## Answers

### 1
Clients are able to send messages to the server, which then flips a coin to see if the message should be forwarded or not. The logic can be seen from the observations above.
### 2
There is a total ordering of messages sent out from the server, as the server maintains a queue and orders the messages to be sent out. So all the clients receive messages in the same order. This doesn't fully mitigate the causality violation possibility all together, as the order of receiving messages at the server is still not the same as the order in which the clients sent out the messages. 

## Part 1, part 3 - Vector clock implementation

This entire section will be dedicated to answering how vector clock has been implemented, and how it is able to detect causality violations, if there are any. 
To simulate a sample scenario, I have included the following send messages:
```go
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
```

Alternatively, you may also enable periodic pinging by uncommenting the line below it:
```go
go enablePeriodicPinging()
```

The periodic pinging proves that the system works with atleast 10 clients. You can increase the number of clients by changing 

### Editable sections
Like the previous part, the editable sections have been marked with 
```go
//START
...
//END
```

### Design
The code has been designed to send and receive messages, just like the previous section. But, instead of a scalar clock, the client maintains a vector clock. Every time there is a send event, the client will increment it's own counter in the vector clock and send it out. 

Everytime there is a receive event, the client will compare every element of it's own clock with every element of the incoming clock. This is where it is able to detect causality violations. If the incoming message has a clock value that is lesser than the clock value for some other client, then it means that there is a potential causality violation. The code to detect causality violation is as follows:

```go
func (client *Client) UpdateVectorClock(V1 []int, V2 []int, from int) []int{
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
		if (V2[i] > V1[i] && i == from){
			// there is a potential causality violation
			clnt.Println("There is a potential causality violation!!")
		}
		V2[i] = max(V2[i], V1[i])
	}
	return V2
}
```
However, due to my excellent engineering skills, my code has been designed in a way that causality violations are impossible! This is because on the sending side the clients send deterministically. Once the server receives a message, it puts the messages into a queue! So the messages are sent in the order they are received. Further the queue is locked when the server is sending messages from it. This ensures that the queue is not modified during the process of sending. So it is impossible for the clients to receive in an order that is different from the order that is sent from the server. 

### Observations and output
In the `~/vector_clock` directory, run the following command:
```shell
go build && ./core
```
The sample scenario given above is being run. 

Upon executing the above command, the clients start sending to the server:

![image](https://github.com/fauzxan/distributed_systems/assets/92146562/bee9515d-0baa-4870-8a16-91ae1fc624d5)

The server receives messages and stores them in a queue. The queue decides weather to drop or forward the message. If it chooses to forward the message, then all the client will start receiving it. 

In the screenshot below, we can see that the server drops the message from 4, and forwards the message from 1:

![image](https://github.com/fauzxan/distributed_systems/assets/92146562/42ce1cdf-1f53-4fce-8b7e-e351840dcb9f)

Once done, the queue state is updated, and all the other node's messages are also received by the queue. Server drops some, and forwards some...

![image](https://github.com/fauzxan/distributed_systems/assets/92146562/b4e24747-0099-4946-8cb4-654257d4b377)

Once the queue state is empty, and there are no more messages to forward, the terminal shows `[Server Ping]` messages. At this point, you can press `ctrl+c` to see the order of messages recieved by the client.

![image](https://github.com/fauzxan/distributed_systems/assets/92146562/5b3b5d5a-0058-40f0-9c4b-28d8c3e9e3ed)

We can see that:

1. The clients don't receive messages from themselves.
2. The order of messages recieved by all the clients are the same.


# Part 2
This part has been completely implemented in the `~/bully_algorithm` folder. So all commands listed below must be issued from there. Open up as many terminals as you want --> Each terminal will become a client when you run `go run main.go` from the parent directory. Go routines are used throughout the code to get a consensus across all the nodes, but the nodes themselves are not go routines, they are complete processes running in more than one terminal. This behaviour is in line with the requirements as confirmed with Professor Sudipta.

Another thing to note here is that, in some of the parts in this part, we are introducing some artificial delay using time.Sleep() to achieve some of the scenarios listed in the parts. We know that TCP orders messages in a queue. So this time.Sleep() tends to mess up the ordering of messages. So if you see the messages printed out of order, please do ignore it. Focus on the logic presented in the explanation. Rest assured, taking away the time.sleep() will print the messages in order as expected.

New nodes joining the network usually obtain a copy of the clientlist by contacting some arbitrary node. To simulate this behaviour, a `clientlist.json` has been included in the code repository. The clientlist is a data source for the newly spun up client to grab hold of an address from the address space. If the `clientlist.json` is is not empty when the first ever node in the network is spun up, you will see a message like this:

![image](https://github.com/fauzxan/distributed_systems/assets/92146562/3d501093-61ac-427d-b354-aa0f232f719a)

However, shortly after it's spun up, the the "dead entry" of node 2 will be removed from the `clientlist.json`, so future nodes joining the network won't have to see the dead entry in the file. As such, the clientlist file is self-correcting. But, if you wish to start with a clean slate of nodes, then feel free to make sure that the `clientlist.json` file has an empty json string: 
```json
{}
```

Upon spinning up any node, you will also be presented with a menu as follows

```shell
Press 1 to enter something into my replica
Press 2 to view my replica
Press 3 to view my clientlist
```

If you are feeling old, please feel free to increase the `time.Sleep()` aka the timeout duration to something longer, so you can see the output clearly 😄. As of submission, it has been set to 1 seconds- so that the messages are sent and received on time. A longer timeout could result in TCP messages not reaching on time, or reaching after some other message. But, you may enter the sleep duration in the following methods to change the timout:

- InvokeReplicaSync() - coordinators detect failed nodes using this
- CoordinatorPing() - non-coordinators detect coordinator failures using this

The time.Sleep() duration in these two methods determines the timeout interval for failure detection, because that's how long the nodes wait for pinging each other. 

## Answers

### 1
In order to see the joint synchronization, open up 2-3 terminal for simplicity sake, and run `go run main.go`. You will see that each client now has a unique id assigned to it, and an IP address. You will also see a menu on each terminal. Press 2 to verify that the replica is empty in each terminal. Once done, on one of the terminals, press 1, and then enter the number you want to enter into the replica of that node. 
Client 0, 1, 2 started in three terminals respectively:
![image](https://github.com/fauzxan/distributed_systems/assets/92146562/86db6bba-f419-47af-8aac-b44a347f949f)

![image](https://github.com/fauzxan/distributed_systems/assets/92146562/b843bfa6-8e2f-44c8-9feb-de86b0c35b2a)

![image](https://github.com/fauzxan/distributed_systems/assets/92146562/14cea6b4-1b09-4c15-98bf-1a09c8ce1f67)

The server will synchronize once every 10 seconds. So you may enter any number of entries into the replicas in any of the terminals, and they will all be synchronized in one or two synchronziation cycles. In the SS below, there were some entries in node 0 and node 1 in the replica. They were synchronized within 2 cycles:

![image](https://github.com/fauzxan/distributed_systems/assets/92146562/0b6f3115-a1e6-460e-924e-66dc68081140)

![image](https://github.com/fauzxan/distributed_systems/assets/92146562/9073d449-4194-4434-8daf-b0f06ce3aaf6)

### 2
In order to measure the best and the worst case times, we make use of a timer to count the time elapsed between when the discover phase started in the clients that detect the election, and when the election ended. 
Four clients were spun up, and a failure was triggerred with the fourth client. The number of clients that detect the failure can be controlled by the controlling the timeout. It is possible that only one node detects the failure if the timeout is large enough. As mentioned, the timeout to detect the coordinator failure needs to be set in `CoordinatorPing()`. However, for the screenshots below, the timeout was only 1s. So multiple nodes detected the failure, and the time between discovery and victory receival (for non victors)/ victory sending(for victor) was calculated for all the nodes.

#### Worst case scenario:
This happens when the client with the lowest id is the one to detect the node failure. As shown below, the lowest id client took almost 3 ms to finish the election. 

![image](https://github.com/fauzxan/distributed_systems/assets/92146562/cbe40bb1-dbd8-4604-9234-332d15d5c1d8)

#### Best case scenario
This happens when the second highest id node is the one to detect the failure. This took around half the time as the highest id node:

![image](https://github.com/fauzxan/distributed_systems/assets/92146562/ca4666ad-4a10-4ade-a6b9-231aadf58ef7)

The node with the middle id took more time than highest id, but less time than lowest id node:

![image](https://github.com/fauzxan/distributed_systems/assets/92146562/6d804cf9-01ce-4437-9e99-2fff39e8e47d)



### 3
We will now consider a case where a node fails during the election. In order to achieve this, we will introduce a small delay when sending out victory messages, so we get sufficient time to fail 
a) the coordinator
b) non-coordinator node

The new coordinator stops two seconds before sending to the other node. This will give us enough time to fail the coordinator/ non-coordinator node. The delay has been introduced as follows:
```go
func (client *Client) Announcevictory(){ // Make yourself coordinator
	send := Message{Type: VICTORY, From: client.Id, Clientlist: client.Clientlist}
	reply := Message{}
	fmt.Println("No higher id node found. I am announcing victory! Current clientlist (you may fail clients now):", client.Clientlist)
	client.Printclients()
	for id, ip := range client.Clientlist{
		if (id == client.Id){continue}
		clnt, err := rpc.Dial("tcp", ip)
		if (err != nil){
			fmt.Println("Communication to", id, "failed")
			delete(client.Clientlist, id)
			go client.updateclientlist()
			continue
		}
		err = clnt.Call("Client.HandleCommunication", send, &reply)
		if (err != nil){
			fmt.Println(err)
			continue
		}
		fmt.Println("Victory message sent to", id, "at", ip)
		if (reply.Type == "ack" && id > client.Id){
			// if you are announcing vicotry and you receive an ack message from a higher id, then just stop announcing
			fmt.Println("Message sent to", id, "successfully")
			fmt.Println("Client", id, "is awake")
			break
		}else if(reply.Type == ACK){
			fmt.Println("Client", id, "acknowledged me as coordinator")
		}
		time.Sleep(2 * time.Second) // <-- DELAY IS BEING INTRODUCED HERE
	}
	go client.InvokeReplicaSync()
}
```

#### Part a
This part will fail the coordinator after discovery, but before announcement. 

For this, we spun up four nodes in four different terminals. (You may do the same with 3, but I just want to demonstrate that this system works with more nodes as well)
In the screenshot given below, we can see that the last node, with id 3 just got created and is trying to announce to everyone that it is the coordinator. We can also see that the discovery phase is done. It is in the process of announcing to everyone. 

![image](https://github.com/fauzxan/distributed_systems/assets/92146562/189eaa16-30c7-47c1-8ff1-114332f5a417)


Above, we see that the new coordinator has already sent messages to node 1 and node 2 who have acknowledge it as the coordinator. However, node 0 has not even received a victory message from node 3. So according to node 0, node 2 is the coordinator. But node 1 and node 2 think 3 is the coordinator. 

##### Node 2:
It was sending out periodic syncs-> then it received victory message from 3-> stopped sending out syncs

Because of the pinging algorithm, where nodes ping the coordinator periodically, this node was able to detect that the coordinator failed, and then started an election. 
![image](https://github.com/fauzxan/distributed_systems/assets/92146562/f565bb67-8958-4837-bc8d-2e1369188cdc)
##### Node 1: 

It was received periodic syncs-> then it received victory message from 3-> stopped sending out ping to 2 -> started sending out ping to 3

Again, the pinging failed, so it started election, and lost to 2. 
![image](https://github.com/fauzxan/distributed_systems/assets/92146562/b6218c18-99ff-47fc-b19a-cadef100afc2)
##### Node 0:
This node never knew any new coordinator other than 2. Hence it was, is, and continues to ping 2 only. In the middle it receives a victory message from 2, but this doesn't matter to 0, as it already believes that 2 is it's coordinator. 

![image](https://github.com/fauzxan/distributed_systems/assets/92146562/0c6e115f-8ad9-4d34-96ea-e3c77401d9e0)

### Part b
In this part we will explore what happens when an arbitrary node fails during the election. As we shall see, the system I have built is self correcting. This means that the failed node will eventually be noticed by the elected coordinator, and the coordinator will take it out of its Clientlist. The other nodes will definitely detect that there is a failed node in the network as soon as they try to communicate with it. This is because all the `rcp.Dial()` instances have a code block that removes the node from their own Clientlist if the node is unable to contact the supposedly failed node. 
```go
clnt, err := rpc.Dial("tcp", coordinator_ip)
if (err != nil){
	fmt.Println("There was an error trying to connect to the coordinator")
	fmt.Println(err)
	fmt.Println("Invoking election")
	delete(client.Clientlist, client.Coordinator_id) // Delete from own clientlist
	client.updateclientlist() 
	go client.Sendcandidacy()
	return
}
err = clnt.Call("Client.HandleCommunication", send, &reply)
if (err != nil){
	fmt.Println("There was an error trying to connect to the coordinator")
	fmt.Println(err)
	fmt.Println("Invoking election")
	delete(client.Clientlist, client.Coordinator_id) // Delete from own clientlist
	go client.Sendcandidacy()
	client.updateclientlist() 
	return
}
```

However, as per requirements of the second part of the assignment, the non-coordinator nodes do not really communicate with each other, due to which the failed node's entry will continue to persist in the Clientlist of non-coordinator nodes. The non-coordinator nodes will remove them from their Clientlist as soon as they get into a situation where they start attempting to communicate with the failed node- i.e., when the non-coordinator node becomes the coordinator, it will try to `InvokeReplicaSync()` with everyone in it's Clientlist.  

In the image below, we see that the client 1 has failed after election started.

![image](https://github.com/fauzxan/distributed_systems/assets/92146562/80797c63-7939-4dcd-bf1a-241dec48d1ee)

And, upon pressing `3` to view clientlist, we see that 1 is indeed not there in its clientlist anymore.

![image](https://github.com/fauzxan/distributed_systems/assets/92146562/74519aa9-3c38-483a-b4df-cbd5c135cc42)


### 4
In order to see a scenario, where multiple clients start the election simultaneously, spin up three clients by running `go run main.go` in three separate clients simultaneously. Then kill the coordinator terminal by pressing `ctrl+c`. The discovery phase in the other two terminals would've begun almost at the same time (separated by a few milliseconds). Previously, in order to slow down the message output, the time.Sleep() value of the CoordinatorPing() method was set to 10 seconds. However, in order to bring about the behaviour required by this part, the time.Sleep() was set to 1 second. 
```go
func (client *Client) Coordinatorping(){ // communicatetocoordinator
	if (len(client.Clientlist) == 1){
		client.Coordinator_id = client.Id
		return
	}
	for {
		time.Sleep(1*time.Second) // <-- this line was modified
		// logic to dial coordinator tcp, and send message to the server.
	}
}
```

As such, we are able to observe from the following screenshots, that the two terminals simultaneously started election at `11:33:45:303`. 

![image](https://github.com/fauzxan/distributed_systems/assets/92146562/445bb0ad-9f2a-4781-ae6a-2d2be2551081)

![image](https://github.com/fauzxan/distributed_systems/assets/92146562/9a10ccfb-d768-461c-8c28-ea772c5a0dd3)


### 5
For an arbitrary node to leave the network, just hit `ctrl+c` in the terminal. 

![image](https://github.com/fauzxan/distributed_systems/assets/92146562/bbe0474f-1830-4c95-ac02-456fda28d8aa)

There are two cases:
1. If the arbitrary node is the coordinator, then the rest of them will figure out that the coordinator is down when trying to ping it. || All clients ping the coordinator every 10 seconds to see if it is alive.![image](https://github.com/fauzxan/distributed_systems/assets/92146562/b83b4b8b-3453-4359-a46e-2699c7fb8f25)

2. If the arbitrary node is not the coordinator, then the coordinator will detect the failure, and update it's clientlist. So next time it sends replica sync request, it won't send to the failed node. The rest of the nodes still have the failed node in their entry. But these nodes will remove the failed nodes as soon as they try and communicate with it. However, as per implementation, the non coordinator nodes never really communicate within themselves, therefore, they don't remove the failed node unless they themselves become the coordinator. ![image](https://github.com/fauzxan/distributed_systems/assets/92146562/6a61ecc7-0a42-49ce-97a8-33a086e8b1eb)
