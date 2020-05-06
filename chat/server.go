package chat

import (
	"encoding/json"
	"fmt"
	ws "gameserver/wsserver"
	"log"
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
	// go LettersIns.Start()
	go Router(ChanForWS)
	// go UsersOnl.pushOnline(ChanForWS)
}

//func AddletForChat(letter ws.Letter) {
//	LettersIns.Add(Letter{letter.ClientID, letter.LetterType, letter.Scroll})
//}

// ReloadLetterFromChatToWS - test
//func ReloadLetterFromChatToWS(ChanForWS chan ws.Letter) { //Letters *ws.Letters
//	ticker := time.Tick(Update * time.Millisecond)
//	for range ticker {
//		if len(LettersFor) > 0 {
//			for _, let := range LettersFor {
//				ChanForWS <- ws.Letter(let)
//				LettersFor.DelFirstL()
//			}
//		}
//	}
//	//newLet := make(ws.Letters, 0, len(*letters))
//	//ticker := time.Tick(Update * time.Millisecond)
//	//for range ticker {
//	//	LetForS := LettersFor.GetAllLetters()
//	//	// fmt.Println("LetForS :", LetForS)
//	//	// fmt.Println("LettersFor :", LettersFor)
//	//	if len(LetForS) > 0 {
//	//		for i := range LetForS {
//	//			newLet = append(newLet, ws.Letter{LetForS[i].ClientID, LetForS[i].LetterType, LetForS[i].Scroll})
//	//		}
//	//		(*letters).PushMore(newLet)
//	//		LettersFor.DelById(len(newLet))
//	//		newLet = make(ws.Letters, 0, len(*letters))
//	//		fmt.Println("  reload mm ")
//	//	}
//	//}
//}

// func UpdateOnlineChat(conn *ws.Connections) {
// 	ticker := time.Tick(SpeedUpdateOnChat * time.Millisecond)
// 	for range ticker {
// 		new := conn.GetOnlineClients()
// 		// fmt.Println("conn: ", conn)
// 		ready := make(UsersOnline, 0, MaxConnections)
// 		for i := range new {
// 			ready = append(ready, UserOnline{new[i].ID, new[i].Nick})
// 		}
// 		UsersOnl.Push(ready)
// 		// fmt.Println("UsersOnl: ", UsersOnl)
// 	}
// }

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
				// fmt.Println("newOnline: ", newOnline)
				// fmt.Println("UsersOnl: ", UsersOnl)
			}

		default:
			// ignore mess
			fmt.Println("Chat, type messages incorrect: ", letter.LetterType)

		}
	}
}
