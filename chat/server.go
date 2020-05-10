package chat

import (
	"encoding/json"
	"fmt"
	"log"
)

var FromChatChan chan []byte = make(chan []byte, 1000)

var InChatChan chan []byte = make(chan []byte, 1000)

func GetFromChatChan() chan []byte {
	return FromChatChan
}

func GetInChatChan() chan []byte {
	return InChatChan
}

// UsersOnl - letters for send to users
var UsersOnl UsersOnline = make(UsersOnline, 0, MaxConnections)

func Start() {

	go Router()
	fmt.Println("init chat")

}

// Router chat logic
func Router() {
	for letter := range InChatChan {

		let := Letter{}
		err := json.Unmarshal(letter, &let)
		if err != nil {
			log.Fatalln(err)
		}

		switch let.LetterType {

		case "2002":
			// mess for all  // AllChatForUser AllChatFromUser
			allU := UsersOnl.GetAllUsersID()
			mes := &AllChatFromUser{}
			err := json.Unmarshal([]byte(let.Scroll), &mes)
			if err != nil {
				log.Println(err)
			}

			message, err := json.Marshal(AllChatForUser{let.ClientID, mes.Text})
			if err != nil {
				log.Printf("error: %v", err)
			}

			for _, us := range allU {
				if us != let.ClientID {
					send, err := json.Marshal(Letter{us, let.LetterType, string(message)})
					if err != nil {
						log.Printf("error: %v", err)
					}
					FromChatChan <- send
				}
			}

		case "2003":
			// privat mess
			mes := &PrivatMessFrom{}

			err := json.Unmarshal([]byte(let.Scroll), &mes)
			if err != nil {
				log.Println(err)
			}
			// PrivatMessFrom PrivatMessFor
			message, err := json.Marshal(PrivatMessFor{let.ClientID, mes.Text})
			if err != nil {
				log.Printf("error: %v", err)
			}

			send, err := json.Marshal(Letter{mes.ForID, let.LetterType, string(message)})
			if err != nil {
				log.Printf("error: %v", err)
			}
			FromChatChan <- send

		case "2005":
			// Search By  Nick
		//	mes := &SerachByNick{}
		//	err := json.Unmarshal([]byte(let.Scroll), &mes)
		//	if err != nil {
		//		log.Println(err)
		//	}
		//	user, res := UsersOnl.GetUserByNick(mes.Nick)
		//	if res {
		//		mes.ID = user.ID
		//		jsonData, err := json.Marshal(mes)
		//		if err != nil {
		//			log.Printf("error: %v", err)
		//			// break
		//		}
		//		send, err := json.Marshal(Letter{let.ClientID, let.LetterType, string(jsonData)})
		//		if err != nil {
		//			log.Printf("error: %v", err)
		//		}
		//		FromChatChan <- send
		//	}
		case "2550":
			//only check
			if let.ClientID == 87654321 {

				newOnline := make(UsersOnline, 0, 500)

				err := json.Unmarshal([]byte(let.Scroll), &newOnline)
				if err != nil {
					fmt.Println(err)
				}
				UsersOnl.Push(newOnline)
				UsersOnl.pushOnlineToClient()

			}

		default:
			// ignore mess
			fmt.Println("Chat, type messages incorrect: ", let.LetterType)

		}
	}
}
