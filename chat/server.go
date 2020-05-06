package chat

import (
	"encoding/json"
	"fmt"
	"log"
	ws "projecta/wsserver"
)

// var FromChatChan chan Letter = make(chan Letter, StackMessages)

var InChatChan chan Letter = make(chan Letter, StackMessages)

//func GetFromChatChan() chan Letter {
//	return FromChatChan
//}

func GetInChatChan() chan Letter {
	return InChatChan
}

// LettersIns - letters inside
var LettersIns Letters = make(Letters, 0, StackMessages)

// LettersFor - letters for send to users
var LettersFor Letters = make(Letters, 0, StackMessages)

// UsersOnl - letters for send to users
var UsersOnl UsersOnline = make(UsersOnline, 0, StackMessages)

func Start(ChanForWS chan ws.Letter) {

	go Router(ChanForWS)
}

// Router chat logic
func Router(ChanForWS chan ws.Letter) {
	for letter := range InChatChan {
		switch letter.LetterType {

		case "2002":
			// mess for all  // AllChatForUser AllChatFromUser
			allU := UsersOnl.GetAllUsersID()
			mes := &AllChatFromUser{}
			err := json.Unmarshal([]byte(letter.Scroll), &mes)
			if err != nil {
				log.Println(err)
			}

			message, err := json.Marshal(AllChatForUser{letter.ClientID, mes.Text})
			if err != nil {
				log.Printf("error: %v", err)
				// break
			}

			for _, us := range allU {
				if us != letter.ClientID {
					ChanForWS <- ws.Letter(Letter{us, letter.LetterType, string(message)})
				}
			}

		case "2003":
			// privat mess
			mes := &PrivatMessFrom{}

			err := json.Unmarshal([]byte(letter.Scroll), &mes)
			if err != nil {
				log.Println(err)
			}
			// PrivatMessFrom PrivatMessFor
			message, err := json.Marshal(PrivatMessFor{letter.ClientID, mes.Text})
			if err != nil {
				log.Printf("error: %v", err)
				// break
			}

			ChanForWS <- ws.Letter(Letter{mes.ForID, letter.LetterType, string(message)})

		case "2005":
			//only check
			mes := &SerachByNick{}
			err := json.Unmarshal([]byte(letter.Scroll), &mes)
			if err != nil {
				log.Println(err)
			}

			user, res := UsersOnl.GetUserByNick(mes.Nick)

			if res {
				mes.ID = user.ID
				jsonData, err := json.Marshal(mes)
				if err != nil {
					log.Printf("error: %v", err)
					// break
				}
				ChanForWS <- ws.Letter(Letter{letter.ClientID, letter.LetterType, string(jsonData)})
			}

		case "2550":
			//only check
			//fmt.Println("letter: ", letter)
			if letter.ClientID == 87654321 {

				newOnline := make(UsersOnline, 0, 500)

				err := json.Unmarshal([]byte(letter.Scroll), &newOnline)
				if err != nil {
					fmt.Println(err)
				}
				UsersOnl.Push(newOnline)
				UsersOnl.pushOnlineToClient(ChanForWS)

			}

		default:
			// ignore mess
			fmt.Println("Chat, type messages incorrect: ", letter.LetterType)

		}
	}
}
