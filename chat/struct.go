package chat

// all chat from user
type AllChatFromUser struct {
	// FromID int    `json:"f"`
	Text string `json:"text"`
}

// AllChatForUser - all chat for user
type AllChatForUser struct {
	FromID int    `json:"id"`
	Text   string `json:"text"`
}

//// StandardMess - mess for all
//type StandardMess struct {
//	FromID int    `json:"fromid"`
//	Text   string `json:"text"`
//}

// WhoOnline - struct for func check online
type WhoOnline struct {
	Users []Who `json:"users"`
}

// Who - struct for WhoOnline
type Who struct {
	ID   int    `json:"id"`
	Nick string `json:"nick"`
}

// ================================  Privat Message

// PrivatMessFrom - for private mess
type PrivatMessFrom struct {
	ForID int `json:"id"`
	// FromID int    `json:"fid"`
	Text string `json:"text"`
}

// PrivatMessFor - for private mess
type PrivatMessFor struct {
	//ForID int `json:"f"`
	FromID int    `json:"id"`
	Text   string `json:"text"`
}

// ================================  SerachByNick
// SerachByNick - for private mess
type SerachByNick struct {
	ID   int    `json:"id"`
	Nick string `json:"nick"`
}
