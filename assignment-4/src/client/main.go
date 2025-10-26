package main

import (
	"bufio"
	"fmt"
	"log"
	"net/rpc"
	"os"

	"github.com/matic113/rpc-demo/src/common"
)

// Client
func main() {
	fmt.Printf("Enter your name: ")
	var Name string
	fmt.Scanln(&Name)
	fmt.Printf("Welcome %s!, You've joined the chat room. Type a message to see chat history.\n", Name)

	// Connect to the RPC server
	rpc_client, error := rpc.Dial("tcp", "127.0.0.1:1234")
	if error != nil {
		log.Fatal("Couldn't connect to server:", error)
		return
	}
	defer rpc_client.Close()

	for {
		fmt.Print("Enter message (or 'exit' to quit): ")
		var message string

		// read user input
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		message = scanner.Text()
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
	fmt.Println()
	fmt.Println("----- Chat History -----")
	for _, msg := range chatHistory {
		fmt.Printf("%s: %s\n", msg.SenderName, msg.Message)
	}
	fmt.Println("------------------------")
	fmt.Println()
}
