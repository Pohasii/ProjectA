package wsserver

// Message - slice bytes
//[]byte
type Message []byte

// Messages struc for queue messages for/from a client
//  type []Template
type Messages []Message

// AddMessage - func AddMessage(message Message)
// add new messages in array
// Message []byte
func (m *Messages) AddMessage(message Message) {
	*m = append(*m, message)
}

// GetMessages - func (ver *Messages) GetMessages() Messages
// get array of messages
// func (m *Messages) GetMessages() Messages {
//	return *m
// }

// DelFirstM - func (ver *Messages) DelFirstM()
// delete first messages in array
func (m *Messages) DelFirstM() {
	if len(*m) >= 2 {
		*m = append((*m)[1:])
	} else {
		*m = make(Messages, 0, cap(*m))
	}
}

// Letter - message type for global array
type Letter struct {
	ClientID   int
	LetterType string
	Scroll     string
}
