package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	ch "projecta/chat"
	ws "projecta/wsserver"
	"strings"
)

func main() {

	// init address for WS Server
	func() {

		var ip string
		var port string
		fmt.Print("Enter IP:")
		fmt.Fscan(os.Stdin, &ip)
		fmt.Print("Enter Port:")
		fmt.Fscan(os.Stdin, &port)
		adr := strings.Join([]string{ip, port}, ":")
		ws.Addr = adr
		fmt.Println("The server started at this address: ", adr)
	}()

	// =====================================================================
	// ws server
	Clients := &ws.Conns

	ChanFromWS := ws.GetFromConnChan()
	ChanForWS := ws.GetOutChan()

	go ws.Start()

	// =====================================================================
	// Chats

	ch.Start() //
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

	//fmt.Println(ws.LettersFrom)

	for let := range ChanFromWS {
		go func() {

			letter := letterType{}
			err := json.Unmarshal(let, &letter)
			if err != nil {
				log.Fatalln("main.go: 65 line: ", err)
			}

			if letter.LetterType == "1001" {

				scrollFrom := &registerNewUs{}
				err := json.Unmarshal([]byte(letter.Scroll), &scrollFrom)
				if err != nil {
					log.Println(err)
				}

				// res := Clients.Add(ws.LettersFrom[i].ClientID, scrollFrom.Nick)
				(*Clients)[letter.ClientID].Nick = scrollFrom.Nick

				jsonData, err := json.Marshal(registerNewUsTrue{letter.ClientID, (*Clients)[letter.ClientID].Status})
				if err != nil {
					log.Printf("error: %v", err)
					// break
				}

				send, err := json.Marshal(letterType{letter.ClientID, "1001", string(jsonData)})
				if err != nil {
					log.Println(err)
				}

				// fmt.Println(send)
				ChanForWS <- send
				(*Clients).PushOnlineClientsToChat()
			}

			switch letter.LetterType[0:1] {

			case "1":
				// auth
				fmt.Println("auth in dev")
			case "2":
				fmt.Print("send to chat: ")
				fmt.Println(letter.Scroll)

				send, err := json.Marshal(letter)
				if err != nil {
					log.Println(err)
				}
				InChatChan <- send // ch.Letter(letter)
			default:
				fmt.Println("incorrect message from userID: ", letter.ClientID)
			}
		}()
	}

}

type registerNewUs struct {
	Nick string `json:"n"`
}

type registerNewUsTrue struct {
	ID     int  `json:"id"`
	Status bool `json:"s"`
}

type letterType struct {
	ClientID   int
	LetterType string
	Scroll     string
}
