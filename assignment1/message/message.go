package message

import "fmt"

type Message struct{
	Counter int
	Message string
	From int
}

func (message *Message) Show(){
	fmt.Println("Counter: ", message.Counter, "\nMessage:", message.Message, "\nFrom:", message.From)
}