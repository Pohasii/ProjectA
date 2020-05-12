package main

import (
	"encoding/json"
	"fmt"
	"log"
	ch "projecta.com/me/chat"
	cl "projecta.com/me/client"
	"projecta.com/me/setenv"
	vr "projecta.com/me/verification"
	ws "projecta.com/me/wsserver"
)

func main() {

	// init address for WS Server
	func() {
		fmt.Println("Starting..!")
		// upload ENV form config.yml
		setenv.SetupConfigToENV()
	}()

	// init services
	func() {
		// ================================= auth
		go vr.Server()

		// ================================= client
		go cl.Start()

		// ================================= ws server
		go ws.Start()

		// ================================= Chats
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

	// client
	ChanFromClient := cl.GetChanFromClient()
	ChanInClient := cl.GetChanInClient()

	// reload message from chat
	go func() {
		for {
			select {
			case let := <-FromChatChan:
				ChanForWS <- let
			case let := <-ChanFromClient:
				letter := letterType{}
				err := json.Unmarshal(let, &letter)
				if err != nil {
					log.Println("main.go: ", err)
				}
				switch letter.LetterType {
				case "2550":
					InChatChan <- let
				default:
					ChanForWS <- let
				}
			}
		}

	}()

	fmt.Println("The server started!")

	// =====================================================================
	// router sms
	for let := range ChanFromWS {
		go func() {

			letter := letterType{}
			err := json.Unmarshal(let, &letter)
			if err != nil {
				log.Println("main.go: ", err)
			}

			fmt.Println("Сообщения из сокета: ", letter)

			switch letter.LetterType[0:1] {

			case "1":
				// client service
				fmt.Print("send to client service: ")
				fmt.Println(letter.Scroll)
				ChanInClient <- ToByte(letter)
			case "2":
				// chat
				fmt.Print("send to chat: ")
				fmt.Println(letter.Scroll)

				InChatChan <- ToByte(letter) // send // ch.Letter(letter)
			default:
				fmt.Println("incorrect message from userID: ", letter.ClientID)
			}
		}()
	}

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
