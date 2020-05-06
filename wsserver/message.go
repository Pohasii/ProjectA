package wsserver

// Message - slice bytes
//[]byte
type Message []byte

// Letter - message type for global array
type Letter struct {
	ClientID   int
	LetterType string
	Scroll     string
}
