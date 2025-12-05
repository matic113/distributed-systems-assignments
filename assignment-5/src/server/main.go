package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"sync"

	"github.com/matic113/rpc-demo/common"
)

// ClientRecord holds info to push messages to a client
type ClientRecord struct {
	ID     string
	Addr   string
	Client *rpc.Client
	Send   chan common.ChatMessage
}

type Chat struct {
	mu      sync.Mutex
	History []common.ChatMessage
	clients map[string]*ClientRecord
}

// helper: add client and start sender goroutine
func (c *Chat) addClient(id, addr string, rpcClient *rpc.Client) *ClientRecord {
	rec := &ClientRecord{ID: id, Addr: addr, Client: rpcClient, Send: make(chan common.ChatMessage, 50)}
	if c.clients == nil {
		c.clients = make(map[string]*ClientRecord)
	}
	c.clients[id] = rec

	// start sender goroutine
	go func(r *ClientRecord) {
		for msg := range r.Send {
			var ack bool
			if err := r.Client.Call("Client.Receive", msg, &ack); err != nil {
				log.Printf("error calling client %s: %v; removing client", r.ID, err)
				// remove client on error
				c.removeClient(r.ID)
				// broadcast leave
				leave := common.ChatMessage{SenderName: "Server", Message: fmt.Sprintf("User %s left", r.ID)}
				c.broadcast(leave, "")
				return
			}
		}
	}(rec)

	return rec
}

func (c *Chat) removeClient(id string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if rec, ok := c.clients[id]; ok {
		close(rec.Send)
		if rec.Client != nil {
			rec.Client.Close()
		}
		delete(c.clients, id)
	}
}

// Register a client for callbacks
func (c *Chat) Register(args common.RegisterArgs, reply *bool) error {
	// Dial client RPC
	clientRPC, err := rpc.Dial("tcp", args.Addr)
	if err != nil {
		return err
	}

	c.mu.Lock()
	// avoid duplicate IDs
	if c.clients == nil {
		c.clients = make(map[string]*ClientRecord)
	}
	if _, exists := c.clients[args.ID]; exists {
		c.mu.Unlock()
		clientRPC.Close()
		*reply = false
		return fmt.Errorf("id already registered")
	}
	rec := c.addClient(args.ID, args.Addr, clientRPC)

	// send history to new client (do not broadcast these to others)
	history := make([]common.ChatMessage, len(c.History))
	copy(history, c.History)
	c.mu.Unlock()

	for _, m := range history {
		select {
		case rec.Send <- m:
		default:
			// drop if buffer full
		}
	}

	// notify others that user joined
	join := common.ChatMessage{SenderName: "Server", Message: fmt.Sprintf("User %s joined", args.ID)}
	c.broadcast(join, args.ID)

	*reply = true
	return nil
}

// Send a message and return true
func (c *Chat) SendMessage(msg common.ChatMessage, reply *bool) error {
	c.mu.Lock()
	c.History = append(c.History, msg)
	c.mu.Unlock()

	// print to server console
	fmt.Printf("%s: %s\n", msg.SenderName, msg.Message)

	// broadcast to all except sender
	c.broadcast(msg, msg.SenderName)

	*reply = true
	return nil
}

// GetHistory returns the current history
func (c *Chat) GetHistory(_ int, reply *[]common.ChatMessage) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	*reply = append([]common.ChatMessage(nil), c.History...)
	return nil
}

// broadcast sends message to all clients except skipID (if provided). Non-blocking send.
func (c *Chat) broadcast(msg common.ChatMessage, skipID string) {
	c.mu.Lock()
	clients := make([]*ClientRecord, 0, len(c.clients))
	for id, rec := range c.clients {
		if id == skipID {
			continue
		}
		clients = append(clients, rec)
	}
	c.mu.Unlock()

	for _, rec := range clients {
		select {
		case rec.Send <- msg:
		default:
			// drop message for slow client
		}
	}
}

// Server
func main() {
	// Listen on port 1234
	listener, err := net.Listen("tcp", "127.0.0.1:1234")
	if err != nil {
		log.Fatal("listen error:", err)
	}
	fmt.Println("Server running on port 1234...")

	var chat Chat

	// Register service -> Publish the method for clients
	rpc.Register(&chat)

	// Accept connections forever
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("accept error:", err)
			continue
		}
		// ServeConn runs the DefaultServer on a single connection
		go rpc.ServeConn(conn)
	}
}
