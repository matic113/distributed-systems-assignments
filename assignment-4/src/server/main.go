package main

import (
	"fmt"
	"net"
	"net/rpc"

	"github.com/matic113/rpc-demo/src/common"
)

type Chat struct {
	History []common.ChatMessage
}

// Send a message and return updated chat history
func (c *Chat) SendMessage(msg common.ChatMessage, reply *bool) error {
	c.History = append(c.History, msg)
	*reply = true
	// print to server console
	fmt.Printf("%s: %s\n", msg.SenderName, msg.Message)
	return nil
}

func (c *Chat) GetHistory(_ int, reply *[]common.ChatMessage) error {
	*reply = c.History
	return nil
}

// Server
func main() {
	// Listen on port 1234
	listener, _ := net.Listen("tcp", "127.0.0.1:1234")
	fmt.Println("Server running on port 1234...")

	var chat Chat

	// Register service -> Publish the method for clients
	rpc.Register(&chat)

	// Accept connections forever
	for {
		conn, _ := listener.Accept()
		// ServeConn runs the DefaultServer on a single connection
		go rpc.ServeConn(conn)
	}
}
