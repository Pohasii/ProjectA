package client

// CheckErr - struct for error
type CheckErr struct {
	Code int
	Text string
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
	text string `json:"text"`
}

type SetNick1002Succ struct {
	text string `json:"text"`
}

// ===================== First Upload
type FirstUpload struct {
	ID        int    `json:"id"`
	Nick      string `json:"nick"`
	Token     string `json:"token"`
	Friends   []int  `json:"friends"`
}