package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"

	"github.com/matic113/rpc-demo/common"
)

// Client RPC receiver
type Client struct{}

// Receive is called by the server to push messages to this client
func (c *Client) Receive(msg common.ChatMessage, reply *bool) error {
	fmt.Printf("\n[Broadcast] %s: %s\nEnter message (or 'exit' to quit): ", msg.SenderName, msg.Message)
	*reply = true
	return nil
}

func main() {
	fmt.Printf("Enter your name: ")
	var Name string
	fmt.Scanln(&Name)
	fmt.Printf("Welcome %s!, You've joined the chat room. Type a message to see chat history.\n", Name)

	// Start local RPC server for receiving broadcasts
	client := &Client{}
	rpc.Register(client)
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal("couldn't start client listener:", err)
	}
	// accept connections for callbacks
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Println("client accept error:", err)
				continue
			}
			go rpc.ServeConn(conn)
		}
	}()

	// Connect to the RPC server
	rpc_client, error := rpc.Dial("tcp", "127.0.0.1:1234")
	if error != nil {
		log.Fatal("Couldn't connect to server:", error)
		return
	}
	defer rpc_client.Close()

	// Register this client with server so it can receive callbacks
	regReply := false
	regArgs := common.RegisterArgs{ID: Name, Addr: listener.Addr().String()}
	if err := rpc_client.Call("Chat.Register", regArgs, &regReply); err != nil || !regReply {
		log.Fatal("registration failed:", err)
		return
	}

	for {
		fmt.Print("Enter message (or 'exit' to quit): ")
		// read user input
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		message := scanner.Text()

		if message == "exit" {
			break
		}
		chatMessage := common.ChatMessage{SenderName: Name, Message: message}

		reply := false
		var err = rpc_client.Call("Chat.SendMessage", chatMessage, &reply)
		if err != nil {
			log.Println("Error sending message:", err)
			continue
		}

		var chatHistory []common.ChatMessage
		err = rpc_client.Call("Chat.GetHistory", 1, &chatHistory)
		if err != nil {
			log.Println("Error getting chat history:", err)
			continue
		}
		PrintChatHistory(chatHistory)
	}
}

func PrintChatHistory(chatHistory []common.ChatMessage) {
	for _, msg := range chatHistory {
		fmt.Printf("%s: %s\n", msg.SenderName, msg.Message)
	}
	fmt.Println("------------------------")
	fmt.Println()
}
