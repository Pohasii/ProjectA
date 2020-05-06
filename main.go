package main

import (
	"encoding/json"
	"fmt"
	"log"
	ch "projecta/chat"
	ws "projecta/wsserver"
	"time"
)

func main() {

	fmt.Println("Start server :)")

	var stub chan interface{} = make(chan interface{}, 100)

	// =====================================================================
	// ws server
	Clients := &ws.Conns
	go ws.Start(stub, stub)
	LettersFor := &ws.LettersFor

	// =====================================================================
	// Chats

	ch.Start()
	go ch.UpdateOnlineChat(Clients)
	go ch.ReloadLetterFromChatToWS(LettersFor)
	// =====================================================================
	// router sms
	tick := time.Tick(10 * time.Millisecond)

	for range tick {

		//fmt.Println(ws.LettersFrom)

		for i := range ws.LettersFrom {

			if ws.LettersFrom[i].LetterType == "1001" {
				scrollFrom := &registerNewUs{}
				err := json.Unmarshal([]byte(ws.LettersFrom[i].Scroll), &scrollFrom)
				if err != nil {
					log.Println(err)
				}

				// res := Clients.Add(ws.LettersFrom[i].ClientID, scrollFrom.Nick)
				(*Clients)[ws.LettersFrom[i].ClientID].Nick = scrollFrom.Nick

				jsonData, err := json.Marshal(registerNewUsTrue{ws.LettersFrom[i].ClientID, (*Clients)[ws.LettersFrom[i].ClientID].Status})
				if err != nil {
					log.Printf("error: %v", err)
					// break
				}

				ws.LettersFor.Add(ws.Letter{ws.LettersFrom[i].ClientID, "1001", string(jsonData)})

				ws.LettersFrom.DelFirstL()
				continue
			}

			switch ws.LettersFrom[i].LetterType[0:1] {

			case "1":
				// auth
				fmt.Println("auth in dev")
				ws.LettersFrom.DelFirstL()
			case "2":
				fmt.Println("send to chat")
				// ws.LettersFrom[i].
				fmt.Println(ws.LettersFrom[i])
				ch.AddletForChat(ws.LettersFrom[i])
				ws.LettersFrom.DelFirstL()
			default:
				fmt.Println("incorrect message from userID: ", ws.LettersFrom[i].ClientID)
				ws.LettersFrom.DelFirstL()
			}
		}

	}
}

type registerNewUs struct {
	Nick string `json:"n"`
}

type registerNewUsTrue struct {
	ID     int  `json:"id"`
	Status bool `json:"s"`
}

// err := json.Unmarshal(Clients[index].InMess[0], &mes)
// 				if err != nil {
// 					log.Println(err)
// 				}
// jsonData, err := json.Marshal(Message{"test2-ready", m})
// 		if err != nil {
// 			log.Printf("error: %v", err)
// 			break
// 		}
