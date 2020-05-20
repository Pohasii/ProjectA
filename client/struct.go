package client

// CheckErr - struct for error
type CheckErr struct {
	Code int
	Text string
}

// ErrorPattern - struct for answer error's
type ErrorPattern struct {
	Error string `json:"error"`
}

// ===================== auth
type Auth1001 struct {
	Token string `json:"token"`
}

type AuthSuccessful struct {
	Status bool `json:"status"`
}

// ===================== set Nick
type SetNick1002 struct {
	Nick string `json:"nick"`
}

type SetNick1002Err struct {
	Text string `json:"text"`
}

type SetNick1002Succ struct {
	Result bool `json:"result"`
}

// ===================== First Upload
type GetProf struct {
	ID        int    `json:"id"`
	Nick      string `json:"nick"`
	Friends   []int  `json:"friends"`
}

// ===================== search by nick

type GetByNick struct {
	Nick string `json:"nick"`
}

type ResponseForSearchByNick struct {
	ID int `json:"id"`
	Nick string `json:"nick"`
}

// ===================== friend request

type RequestFriend struct {
	ID int `json:"id"`
	Nick string `json:"nick"`
	Response bool `json:"response"`
}

type FromFriendRequest struct {
	ID int `json:"id"`
	Request bool `json:"request"`
}

type removeFriend struct {
	ID int `json:"id"`
}

type friends struct {
	Friends []int `json:"friends"`
}