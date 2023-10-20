package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"core/client"
)

func SetupCloseHandler(server *client.Server) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\n\n\n\rCtrl+C pressed in Terminal!")
		fmt.Println("Closing server channels...")
		server.CloseChannels()
		os.Exit(0)
	}()
}