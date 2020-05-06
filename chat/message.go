package chat

// Message - sdf
// letter string
// from   int
type Message struct {
	letter string
}

// Messages - array Message
type Messages []Message

// AddMessage - func AddMessage(message Message)
func (m *Messages) AddMessage(message Message) {
	*m = append((*m)[:], message)
}

// GetMessages - func (ver *Messages) GetMessages() Messages
// get array of messages
func (m *Messages) GetMessages() Messages {
	return *m
}

// DelFirstM - func (ver *Messages) DelFirstM()
// delete first messages in array
func (m *Messages) DelFirstM() {
	if len(*m) >= 2 {
		*m = append((*m)[1:])
	} else {
		*m = make(Messages, 0, cap(*m))
	}
}
