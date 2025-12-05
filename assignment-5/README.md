# Assignment Text
*Modify the RPC chat system (multiple clients/servers, where requests return full history) to use real-time broadcasting with Go concurrency.
When a client joins, notify all other clients: "User [ID] joined".
When any client sends a message, broadcast it to all other clients (no self-echo).
Use goroutines/channels for concurrent send/receive; sync shared client list with Mutex.
Create a new GitHub repo and attach the link to the assignment. Don't use the submitted repo.*

# Improvements Over Previous Task
This repo improves on the previous task by adding a real-time chat system using Go's goroutines and channels for concurrency. The system allows multiple clients to connect to a server, send messages, and receive broadcasts of messages from other clients in real-time.
