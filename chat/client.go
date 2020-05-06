package chat

import (
	"encoding/json"
	"log"
	ws "projecta/wsserver"
)

type UserOnline struct {
	ID   int    `json:"id"`
	Nick string `json:"n"`
}

type UsersOnline []UserOnline

func (u *UsersOnline) Push(new UsersOnline) {
	*u = new
}

// GetLink - back link to *Letters
func (u *UsersOnline) GetLink() *UsersOnline {
	return u
}

func (u *UsersOnline) GetUserByNick(nick string) (UserOnline, bool) {
	for i, us := range *u {
		if us.Nick == nick {
			return (*u)[i], true
		}
	}
	return UserOnline{0, ""}, false
}

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

func (u *UsersOnline) pushOnlineToClient(ChanForWS chan ws.Letter) {

	jsonData, err := json.Marshal(*u)
	if err != nil {
		log.Printf("error: %v", err)
		// break
	}
	for _, val := range *u {
		ChanForWS <- ws.Letter(Letter{val.ID, "2004", string(jsonData)})
	}
}

// func (u *UsersOnline) pushOnlineToClient(ChanForWS chan ws.Letter) {
// 	push := time.Tick(2000 * time.Millisecond)
// 	old := *u
// 	for range push {
// 		if !reflect.DeepEqual(*u, old) {
// 			jsonData, err := json.Marshal(*u)
// 			if err != nil {
// 				log.Printf("error: %v", err)
// 				// break
// 			}
// 			for _, val := range *u {
// 				ChanForWS <- ws.Letter(Letter{val.ID, "2004", string(jsonData)})
// 			}
// 			old = *u
// 		}
// 	}
// }
