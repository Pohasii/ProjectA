package wsserver

import "time"

// Letter - message type for global array
type Letter struct {
	ClientID   int
	LetterType string
	Scroll     string
}

// Letters - array Letters
type Letters []Letter

// Add (letter Letter) - add new letter to array
// type Letter struct {
//	ClientID   int
//	letterType string
//	Scroll     string
// }
func (l *Letters) Add(letter Letter) {
	// l[len(*l)] = letter
	*l = append((*l)[:], letter)
}

// PushMore Push more letters
func (l *Letters) PushMore(letter []Letter) {
	// l[len(*l)] = letter
	for i := range letter {
		*l = append((*l)[:], letter[i])
	}
}

// GetLink - back link to *Letters
func (l *Letters) GetLink() *Letters {
	return l
}

// DelFirstL - func (ver *Messages) DelFirstM()
// delete first messages in array
func (l *Letters) DelFirstL() {
	if len(*l) >= 2 {
		*l = append((*l)[1:])
	} else {
		*l = make(Letters, 0, cap(*l))
	}
}

// addFor - func for letter sorting for clients
func (l *Letters) addFor() {
	tick := time.Tick(LettersSort * time.Millisecond)
	for range tick {
		if len(*l) > 0 {
			for _, val := range *l {
				for i := range Conns {
					if Conns[i].ID == val.ClientID {
						Conns[i].OutMess.AddMessage([]byte(val.LetterType + val.Scroll))
						(*l).DelFirstL()
					}
				}
			}
		}
	}
}
