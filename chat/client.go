package chat

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
