package verification

// Profile - user profile
type Profile struct {
	ID        int    `json:"id"`
	Nick      string `json:"nick"`
	Password  string
	Login     string
	Token     string `json:"token"`
	Friends   []int  `json:"friends"`
	IsActive  bool
	FirstConn bool
}

type Credentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	ConnType string `json:"conntype"`
}

type Failed struct {
	Error string `json:"error"`
}

type Successful struct {
	Token string `json:"token"`
}

type LastID struct {
	ID    string `json:"id"`
	Value int    `json:"value"`
}
