package main

import (
	"bufio"
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
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter IP:")
		ip, _ := reader.ReadString('\n')
		// convert CRLF to LF
		ip = strings.Replace(ip, "\n", "", -1)

		fmt.Print("Enter Port:")
		port, _ := reader.ReadString('\n')
		// convert CRLF to LF
		port = strings.Replace(port, "\n", "", -1)

		fmt.Println("address:", ip+":"+port)
		ws.Addr = ip + ":" + port
	}()

	fmt.Println("Start server :)")
	// =====================================================================
	// ws server
	Clients := &ws.Conns

	ChanFromWS := ws.GetInChan()
	ChanForWS := ws.GetOutChan()

	go ws.Start(ChanFromWS, ChanForWS)
	// LettersFor := &ws.LettersFor

	// =====================================================================
	// Chats
	InChatChan := ch.GetInChatChan()
	ch.Start(ChanForWS)
	// go ch.UpdateOnlineChat(Clients)
	// go ch.ReloadLetterFromChatToWS(ChanForWS)

	// GetFromChatChan
	// =====================================================================
	// router sms

	//fmt.Println(ws.LettersFrom)

	for letter := range ChanFromWS {

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

			ChanForWS <- ws.Letter{letter.ClientID, "1001", string(jsonData)}
			// fmt.Println((*Clients)[letter.ClientID])
			(*Clients).PushOnlineClientsToChat()
			continue
		}
		switch letter.LetterType[0:1] {

		case "1":
			// auth
			fmt.Println("auth in dev")
		case "2":
			fmt.Print("send to chat: ")
			// ws.LettersFrom[i].
			fmt.Println(letter.Scroll)
			// ch.AddletForChat(letter)
			InChatChan <- ch.Letter(letter)
		default:
			fmt.Println("incorrect message from userID: ", letter.ClientID)
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

// tick := time.Tick(10 * time.Millisecond)
// for range tick {
// 	//fmt.Println(ws.LettersFrom)
// 	for i := range ws.LettersFrom {
// 		if ws.LettersFrom[i].LetterType == "1001" {
// 			scrollFrom := &registerNewUs{}
// 			err := json.Unmarshal([]byte(ws.LettersFrom[i].Scroll), &scrollFrom)
// 			if err != nil {
// 				log.Println(err)
// 			}
// 			// res := Clients.Add(ws.LettersFrom[i].ClientID, scrollFrom.Nick)
// 			(*Clients)[ws.LettersFrom[i].ClientID].Nick = scrollFrom.Nick
// 			jsonData, err := json.Marshal(registerNewUsTrue{ws.LettersFrom[i].ClientID, (*Clients)[ws.LettersFrom[i].ClientID].Status})
// 			if err != nil {
// 				log.Printf("error: %v", err)
// 				// break
// 			}
// 			ws.LettersFor.Add(ws.Letter{ws.LettersFrom[i].ClientID, "1001", string(jsonData)})
// 			ws.LettersFrom.DelFirstL()
// 			continue
// 		}
// 		switch ws.LettersFrom[i].LetterType[0:1] {
// 		case "1":
// 			// auth
// 			fmt.Println("auth in dev")
// 			ws.LettersFrom.DelFirstL()
// 		case "2":
// 			fmt.Println("send to chat")
// 			// ws.LettersFrom[i].
// 			fmt.Println(ws.LettersFrom[i])
// 			ch.AddletForChat(ws.LettersFrom[i])
// 			ws.LettersFrom.DelFirstL()
// 		default:
// 			fmt.Println("incorrect message from userID: ", ws.LettersFrom[i].ClientID)
// 			ws.LettersFrom.DelFirstL()
// 		}
// 	}
// }
