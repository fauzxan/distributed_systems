package client

import (
	"sync"
	"fmt"
)


type Clientlist struct{
	Lock sync.Mutex
	Clients map[int]*Client
} 


func (clientlist *Clientlist) Add(client *Client){
	clientlist.Lock.Lock()
	if clientlist.Clients == nil {
        clientlist.Clients = make(map[int]*Client)
    }
    clientlist.Clients[client.Clientid] = client
	clientlist.Lock.Unlock()
}

func (clientlist *Clientlist) Check(clientid int) bool{
	//check to see if clientid exists within the client map
	clientlist.Lock.Lock()
	defer clientlist.Lock.Unlock()
	_, exists := clientlist.Clients[clientid]
    return exists
}

func (clientlist *Clientlist) Getmax() int{
	// return the maximum key in the client mapping
	clientlist.Lock.Lock()
	defer clientlist.Lock.Unlock()
	max := 0

    for clientid := range clientlist.Clients {
        if clientid > max {
            max = clientid
        }
    }
    return max
}

func (clientlist *Clientlist) PrintClients() {
    clientlist.Lock.Lock()
    

    fmt.Println("Client List:")
    for clientid, client := range clientlist.Clients {
        fmt.Println("Client ID:", clientid, "with counter value", client.Counter)
    }
	clientlist.Lock.Unlock()
}

func (clientlist *Clientlist) CloseChannels(){
	for clientid, client := range clientlist.Clients{
		close(client.Channel)
		fmt.Println("Closed channel for", clientid)
	}
}