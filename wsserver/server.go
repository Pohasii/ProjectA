package wsserver

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	// "crypto/tls"
	"github.com/gorilla/websocket"
)

var OutChan chan []byte = make(chan []byte, 1000)

// var InChan chan Letter = make(chan Letter, 500)
var FromConnChan chan []byte = make(chan []byte, 1000)

func GetOutChan() chan []byte {
	return OutChan
}

func GetFromConnChan() chan []byte {
	return FromConnChan
}

// Conns - all connection clients
var Conns Connections = make(Connections, 0, MaxConnections)

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
func Start() {

	go Conns.CleanOffConn()
	go router()

	flag.Parse()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		ServeWs(w, r, &Conns)
	})

	fmt.Println("The websocket's server started at the ", os.Getenv("WebsocketIP")+":"+os.Getenv("WebsocketPORT"))

	// err := srv.ListenAndServeTLS(*addr,"server.crt", "server.key", nil)
	err := http.ListenAndServe(*addr, nil) // *addr
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}

func router () {
	for let := range OutChan {
		letter := Letter{}
		err := json.Unmarshal(let, &letter)
		if err != nil {
			log.Fatalln(err)
		}

		switch letter.LetterType {
		case "1901":
			if Conns[letter.ClientID].Status == false {
				Conns.DelByID(letter.ClientID)
				fmt.Println("Remove connection:", letter.ClientID)
			}
		case "1902":
			Conns[letter.ClientID].Auth = true
		default:
			for _, conn := range Conns {
				if conn.ID == letter.ClientID && conn.Status != false {
					conn.Send <- []byte(letter.LetterType + letter.Scroll)
				}
			}
		}
	}
}