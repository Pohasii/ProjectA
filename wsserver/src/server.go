package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
)

var OutChan = make(chan []byte, 1000)

// var InChan chan Letter = make(chan Letter, 500)
var FromConnChan = make(chan []byte, 1000)

// Conns - all connection clients
var Conns = make(Connections, 0, MaxConnections)

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
func ServeWs(w http.ResponseWriter, r *http.Request) {
	for {

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("ServeWs, upgrader.Upgrade: ", err)
			return
		}

		fmt.Println("new client: ", conn.RemoteAddr())
		// add conn to

		res := Hub.AddClient(InitClient(conn))
		if res {
			fmt.Println("successful")
		} else {
			fmt.Println("failed to add client")
		}

		//(*Conns).Add(conn)
		//(*Conns)[len(*Conns)-1].start()

	}
}

// Start func Start(client *Client)
// start http Websocket server
func Start() {

	// redis init connections
	SetEnv()
	go router()

	flag.Parse()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		ServeWs(w, r)
	})

	var addr = flag.String("addr", os.Getenv("WebsocketIP")+":"+os.Getenv("WebsocketPORT"), "http service address")
	fmt.Println("The websocket's server started at the ", os.Getenv("WebsocketIP")+":"+os.Getenv("WebsocketPORT"))

	// err := http.ListenAndServeTLS(*addr,"private.crt", "private.key", nil)
	// http.ListenAndServeTLS()

	err := http.ListenAndServe(*addr, nil) // *addr
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}

func router() {
	for let := range OutChan {
		letter := Letter{}
		err := json.Unmarshal(let, &letter)
		if err != nil {
			log.Println(err)
		}

		switch letter.LetterType {
		case "1901":
			if Conns[letter.ClientID].Status == false {
				Conns.DelByID(letter.ClientID)
				fmt.Println("Remove connection: ", letter.ClientID)
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
