package wsserver

import (
	"flag"
	"os"
	"time"
)

// Addr - address for open listening
// using in:
// server/Start()
// main.go
// var Addr string

// var addr = flag.String("addr", Addr, "http service address")
var addr = flag.String("addr", os.Getenv("WebsocketIP")+":"+os.Getenv("WebsocketPORT"), "http service address")

const (
	// ClearPer - Сhecking for messages to send
	// using in clients/CleanOffConn()
	ClearPer = 2000

	// СheckingMessages - Сhecking for messages to send
	// useing in client/writePump()
	СheckingMessages = 25

	// TimeForAuth - Time  For Auth
	// useing in connections/func (c *Connections) Add
	TimeForAuth = 10000


	// MaxConnections - max connections clients
	// using in:
	// clients/CleanOffConn()
	// server/Connections
	MaxConnections = 500

	// MessagesCapacity - Message array Capacity
	// using in:
	// clients/Add
	MessagesCapacity = 200

	// LettersSort - letter sorting frequency
	// using:
	// letters/addFor()
	LettersSort = 10

	// Maximum message size allowed from peer.
	// using in:
	// clients/Add
	maxMessageSize = 512

	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)
