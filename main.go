package main

import (
	"encoding/json"
	"fmt"
	"log"
	ch "projecta/chat"
	cl "projecta/client"
	vr "projecta/verification"
	ws "projecta/wsserver"
)

func main() {

	// init address for WS Server
	func() {

		// var ip string
		// var port string
		// fmt.Print("Enter IP:")
		// fmt.Fscan(os.Stdin, &ip)
		// fmt.Print("Enter Port:")
		// fmt.Fscan(os.Stdin, &port)
		// adr := strings.Join([]string{ip, port}, ":")
		ws.Addr = "192.168.0.65:55443" // adr
		// fmt.Println("The server started at this address: ", adr)
	}()

	// init services
	func() {
		// ================================= auth
		go vr.Server()
		// cl.Start()

		// ================================= client
		go cl.Start()

		// =====================================================================
		// ws server
		go ws.Start()

		// =====================================================================
		// Chats
		ch.Start()
	}()

	// ===================================================================
	// init channels

	// websocket
	ChanFromWS := ws.GetFromConnChan()
	ChanForWS := ws.GetOutChan()

	// chat
	InChatChan := ch.GetInChatChan()
	FromChatChan := ch.GetFromChatChan()

	// reload message from chat
	go func() {
		for val := range FromChatChan {
			ChanForWS <- val
		}
	}()

	// =====================================================================
	// router sms

	for let := range ChanFromWS {
		go func() {

			letter := letterType{}
			err := json.Unmarshal(let, &letter)
			if err != nil {
				log.Fatalln("main.go: ", err)
			}

			switch letter.LetterType[0:1] {

			case "1":
				// auth
				fmt.Print("send to auth: ")
				fmt.Println(letter.Scroll)

			case "2":
				// chat
				fmt.Print("send to chat: ")
				fmt.Println(letter.Scroll)

				// send, err := json.Marshal(letter)
				// if err != nil {
				// 	log.Println(err)
				// }

				InChatChan <- ToByte(letter) // send // ch.Letter(letter)
			default:
				fmt.Println("incorrect message from userID: ", letter.ClientID)
			}
		}()
	}

}

type registerNewUs struct {
	Nick string `json:"nick"`
}

type registerNewUsTrue struct {
	ID int `json:"id"`
	// Status bool `json:"status"`
}

type letterType struct {
	ClientID   int
	LetterType string
	Scroll     string
}

func ToByte(letter letterType) []byte {
	send, err := json.Marshal(letter)
	if err != nil {
		log.Println(err)
	}
	return send
}
