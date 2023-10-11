# 1.1 and 1.2

In order to run this file, go to the parent directory, run `go build`. This will generate a file called `core`. You can run this file as an executable - `.\core`.

# Design:

The server has 5 channels, to enable concurrency. Each time it receievs a message from some client, it will add it to a server queue. This is the list of messages that it needs to send out in FIFO order. 
THe queue will randomly drop messages:

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

## Verification:
Running the code using `.\core` will generate logs that you can use to verify the total ordering of the messages. All in all, there are two things done to ensure that the messages received are in the same order that they were sent:

1. Clients order the messages internally before sending them.
2. Server will only send the messages in the order that it receives them. 