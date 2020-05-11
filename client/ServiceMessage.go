package client

// Letter - message type in system
type Letter struct {
	ClientID   int
	LetterType string
	Scroll     string
}

// Online struct for sent online
type Online struct {
	ConnID int    `json:"connid"`
	UserID int    `json:"userid"`
	Nick   string `json:"nick"`
}
