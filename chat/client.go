package chat

import (
	"encoding/json"
	"log"
)

type UserOnline struct {
	ConnID int    `json:"connid"`
	UserID int    `json:"userid"`
	Nick   string `json:"nick"`
}

type UsersOnline []UserOnline
// type UsersOnline map[int]UserOnline

func (u *UsersOnline) Push(Us UserOnline) {
	add := true
	for _, val := range *u {
		if val.UserID == Us.UserID {
			(*u).DelByID(Us.UserID)
			add = false
			break
		}
	}
	if add {
		*u = append(*u, Us )
	}
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
	return UserOnline{0, 0, ""}, false
}

func (u *UsersOnline) GetAllUsers() UsersOnline {
	return *u
}

func (u *UsersOnline) GetAllUsersID() []int {
	IDs := make([]int, 0, len(*u))
	for _, us := range *u {
		IDs = append(IDs, us.ConnID)
	}
	return IDs
}

func (u *UsersOnline) pushOnlineToClient() {

	type online struct{
		ID int `json:"id"`
		Nick string `json:"nick"`
	}

	onl := []online{}
	for _, val := range *u {
		onl = append(onl, online{val.UserID, val.Nick})
	}

	jsonData, err := json.Marshal(onl)
	if err != nil {
		log.Printf("error: %v", err)
	}
	for _, val := range *u {
		send, err := json.Marshal(Letter{val.ConnID, "2004", string(jsonData)})
		if err != nil {
			log.Printf("error: %v", err)
		}
		FromChatChan <- send
	}
}

// DelByID - delete client from Clients array by id
// expected id (int)
func (u *UsersOnline) DelByID(id int) {

	for i, val := range *u {
		if val.UserID == id {
			switch i {
			case 0:
				*u = append((*u)[1:])
			case len(*u)-1:
				*u = append((*u)[0:id-1])
			default:
				*u = append((*u)[:id], (*u)[id+1:]...)
			}

			break
		}
	}

}
