package chat

// all chat from user
type AllChatFromUser struct {
	// FromID int    `json:"f"`
	Text string `json:"t"`
}

// AllChatForUser - all chat for user
type AllChatForUser struct {
	FromID int    `json:"f"`
	Text   string `json:"t"`
}

// StandardMess - mess for all
type StandardMess struct {
	FromID int    `json:"f"`
	Text   string `json:"t"`
}

// WhoOnline - struct for func check online
type WhoOnline struct {
	Users []Who `json:"u"`
}

// Who - struct for WhoOnline
type Who struct {
	ID   int    `json:"id"`
	Nick string `json:"n"`
}

// PrivatMessFrom - for private mess
type PrivatMessFrom struct {
	ForID int `json:"f"`
	// FromID int    `json:"fid"`
	Text string `json:"t"`
}

// PrivatMessFor - for private mess
type PrivatMessFor struct {
	//ForID int `json:"f"`
	FromID int    `json:"fid"`
	Text   string `json:"t"`
}

// SerachByNick - for private mess
type SerachByNick struct {
	ID   int    `json:"id"`
	Nick string `json:"n"`
}
