package chat

import (
	"encoding/json"
	"log"
)

type UserOnline struct {
	ID   int    `json:"id"`
}

type UsersOnline []UserOnline
// type UsersOnline map[int]UserOnline

func (u *UsersOnline) Push(new UsersOnline) {
	*u = new
}

// GetLink - back link to *Letters
func (u *UsersOnline) GetLink() *UsersOnline {
	return u
}

//func (u *UsersOnline) GetUserByNick(nick string) (UserOnline, bool) {
//	for i, us := range *u {
//		if us.Nick == nick {
//			return (*u)[i], true
//		}
//	}
//	return UserOnline{0, ""}, false
//}

func (u *UsersOnline) GetAllUsers() UsersOnline {
	return *u
}

func (u *UsersOnline) GetAllUsersID() []int {
	IDs := make([]int, 0, len(*u))
	for _, us := range *u {
		IDs = append(IDs, us.ID)
	}
	return IDs
}

func (u *UsersOnline) pushOnlineToClient() {

	jsonData, err := json.Marshal(*u)
	if err != nil {
		log.Printf("error: %v", err)
	}
	for _, val := range *u {
		send, err := json.Marshal(Letter{val.ID, "2004", string(jsonData)})
		if err != nil {
			log.Printf("error: %v", err)
		}
		FromChatChan <- send
	}
}
