package chat

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

// Letter - message type for global array
type Letter struct {
	ClientID   int
	LetterType string
	Scroll     string
}

// Letters - array Letters
type Letters []Letter

// Add (letter Letter) - add new letter to array
// type Letter struct {
//	ClientID   int
//	letterType string
//	Scroll     string
// }
func (l *Letters) Add(letter Letter) {
	// l[len(*l)] = letter
	*l = append((*l)[:], letter)
}

// GetAllLetters - gh
func (l *Letters) GetAllLetters() Letters {
	// l[len(*l)] = letter
	return *l
}

// GetAllLetters - gh
func (l *Letters) DelById(ID int) {
	// l[len(*l)] = letter
	if (ID == 0 || ID == 1) && len(*l) > 0 {
		*l = (*l)[ID:]
	} else {
		*l = (*l)[ID-1:]
	}
}

// GetLink - back link to *Letters
func (l *Letters) GetLink() *Letters {
	return l
}

// DelFirstL - func (ver *Messages) DelFirstM()
// delete first messages in array
func (l *Letters) DelFirstL() {
	if len(*l) >= 2 {
		*l = append((*l)[1:])
	} else {
		*l = make(Letters, 0, cap(*l))
	}
}

func (l *Letters) Start() {

	ticker := time.Tick(Update * time.Millisecond)

	for range ticker {

		for _, let := range *l {

			switch let.LetterType {

			case "2002":
				// mess for all
				allU := UsersOnl.GetAllUsersID()
				mes := &StandardMess{}
				err := json.Unmarshal([]byte(let.Scroll), &mes)
				if err != nil {
					log.Println(err)
				}

				jsonData, err := json.Marshal(StandardMess{mes.FromID, mes.Text})
				if err != nil {
					log.Printf("error: %v", err)
					// break
				}

				for _, us := range allU {
					if us != let.ClientID {
						LettersFor.Add(Letter{us, let.LetterType, string(jsonData)})
					}
				}

				(*l).DelFirstL()

			case "2003":
				// privat mess
				mes := &PrivatMess{}

				err := json.Unmarshal([]byte(let.Scroll), &mes)
				if err != nil {
					log.Println(err)
				}

				jsonData, err := json.Marshal(StandardMess{mes.FromID, mes.Text})
				if err != nil {
					log.Printf("error: %v", err)
					// break
				}

				LettersFor.Add(Letter{mes.ForID, let.LetterType, string(jsonData)})
				(*l).DelFirstL()

			case "2004":
				//only check
				allU := UsersOnl.GetAllUsers()

				jsonData, err := json.Marshal(allU)
				if err != nil {
					log.Printf("error: %v", err)
					// break
				}
				LettersFor.Add(Letter{let.ClientID, let.LetterType, string(jsonData)})
				(*l).DelFirstL()

			case "2005":
				//only check
				mes := &SerachByNick{}
				err := json.Unmarshal([]byte(let.Scroll), &mes)
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
					LettersFor.Add(Letter{let.ClientID, let.LetterType, string(jsonData)})
				}

				(*l).DelFirstL()

			default:
				// ignore mess
				fmt.Println("Chat, type messages incorrect: ", let.LetterType)
				(*l).DelFirstL()
			}
		}
	}
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
