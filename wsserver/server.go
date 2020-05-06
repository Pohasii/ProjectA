package wsserver

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var OutChan chan Letter = make(chan Letter, 500)

var InChan chan Letter = make(chan Letter, 500)

func GetOutChan() chan Letter {
	return OutChan
}

func GetInChan() chan Letter {
	return InChan
}

// Conns - all connection clients
var Conns Connections = make(Connections, 0, MaxConnections)

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
func Start(ChanFromWS, ChanForWS chan Letter) {

	go Conns.CleanOffConn()
	go sortingForUsers()

	flag.Parse()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ServeWs(w, r, &Conns)
	})

	var addr = flag.String("addr", Addr, "http service address")
	err := http.ListenAndServe(*addr, nil) // *addr
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// sortingForUsers - func for letter sorting for clients
func sortingForUsers() { // addFor
	for let := range OutChan {
		for _, conn := range Conns {
			if conn.ID == let.ClientID && conn.Status != false {
				conn.Send <- []byte(let.LetterType + let.Scroll)
			}
		}
	}
}
