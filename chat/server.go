package chat

import (
	"fmt"
	ws "projecta/wsserver"
	"time"
)

// LettersIns - letters inside
var LettersIns Letters = make(Letters, 0, StackMessages)

// LettersFor - letters for send to users
var LettersFor Letters = make(Letters, 0, StackMessages)

// UsersOnl - letters for send to users
var UsersOnl UsersOnline = make(UsersOnline, 0, StackMessages)

func Start() {
	go LettersIns.Start()
}

func AddletForChat(letter ws.Letter) {
	LettersIns.Add(Letter{letter.ClientID, letter.LetterType, letter.Scroll})
}

func ReloadLetterFromChatToWS(letters *ws.Letters) {

	newLet := make(ws.Letters, 0, len(*letters))

	ticker := time.Tick(Update * time.Millisecond)

	for range ticker {
		LetForS := LettersFor.GetAllLetters()
		// fmt.Println("LetForS :", LetForS)
		// fmt.Println("LettersFor :", LettersFor)
		if len(LetForS) > 0 {
			for i := range LetForS {
				newLet = append(newLet, ws.Letter{LetForS[i].ClientID, LetForS[i].LetterType, LetForS[i].Scroll})
			}
			(*letters).PushMore(newLet)
			LettersFor.DelById(len(newLet))
			newLet = make(ws.Letters, 0, len(*letters))

			fmt.Println("  reload mm ")
		}
	}
}

func UpdateOnlineChat(conn *ws.Connections) {
	ticker := time.Tick(SpeedUpdateOnChat * time.Millisecond)

	for range ticker {
		new := conn.GetOnlineClients()
		// fmt.Println("conn: ", conn)
		ready := make(UsersOnline, 0, MaxConnections)
		for i := range new {
			ready = append(ready, UserOnline{new[i].ID, new[i].Nick})
		}
		UsersOnl.Push(ready)
		// fmt.Println("UsersOnl: ", UsersOnl)
	}
}
