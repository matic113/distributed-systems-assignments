package common

type ChatMessage struct {
	SenderName, Message string
}

// RegisterArgs is used by a client to register its callback address with the server
type RegisterArgs struct {
	ID   string
	Addr string
}
