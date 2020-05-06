package wsserver

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// Conns - all connection clients
var Conns Connections = make(Connections, 0, MaxConnections)

// LettersFrom - Letters from all users for router
var LettersFrom Letters

// LettersFor - Letters for all users to send
var LettersFor Letters

// RunCleanConn - defoult false
var RunCleanConn bool = false

// 192.168.0.65
var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// ServeWs handles websocket requests from the peer.
func ServeWs(w http.ResponseWriter, r *http.Request, Conns *Connections) {
	for {

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Where: ", err)
			return
		}

		fmt.Println("new client: ", conn.RemoteAddr())
		// add conn to
		(*Conns).Add(conn)
		(*Conns)[len(*Conns)-1].start()

	}
}

// Start func Start(client *Client)
// start http Websocket server
func Start(sendMessage, incomingMessages chan interface{}) {

	go Conns.CleanOffConn()
	go LettersFor.addFor()

	flag.Parse()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ServeWs(w, r, &Conns)
	})
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
