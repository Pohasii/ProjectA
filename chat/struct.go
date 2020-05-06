package chat

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

// PrivatMess - for private mess
type PrivatMess struct {
	ForID  int    `json:"fod"`
	FromID int    `json:"fid"`
	Text   string `json:"t"`
}

// SerachByNick - for private mess
type SerachByNick struct {
	ID   int    `json:"id"`
	Nick string `json:"n"`
}
