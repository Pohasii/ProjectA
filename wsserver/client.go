package wsserver

import (
	"bytes"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// Client - struct
type Client struct {
	// The websocket connection.
	Conn *websocket.Conn

	// Buffered channel of outbound messages.
	Send    chan []byte
	OutMess Messages
	// InMess  Messages
	ID     int
	Status bool
	Remove bool
	Nick   string
}

func (c *Client) writePump() {

	ticker := time.NewTicker(pingPeriod)
	// CheckTicker := time.Tick(СheckingMessages * time.Millisecond)

	defer func() {
		ticker.Stop()
		c.Conn.Close()
		c.Status = false
		close(c.Send)
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// case <-CheckTicker:
// 	if !c.Status {
// 		return
// 	}
// 	if len(c.OutMess) > 0 {
// 		go func() {
// 			Message := c.OutMess
// 			var count int
// 			for i, mes := range Message {
// 				c.Send <- mes
// 				count = i
// 			}
// 			// c.Send <- c.OutMess.GetMessages()[0]
// 			if count == 0 {
// 				(*c).OutMess = c.OutMess[1:]
// 			} else {
// 				(*c).OutMess = c.OutMess[count:]
// 			}
// 		}()
// 	}

// readPump pumps messages from the websocket connection to the hub.
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {

	defer func() {
		c.Conn.Close()
		c.Status = false
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))

		lettToStr := string(message)
		typeLett := lettToStr[0:4]
		letter := lettToStr[4:]

		// LettersFrom.Add(Letter{c.ID, typeLett, letter})
		InChan <- Letter{c.ID, typeLett, letter}
		// fmt.Println("New message: ", message, " to string: ", string(message))
		// c.Out <- message
		// c.InMess.AddMessage(message)
	}
}

// outgoingMessagesRouter func for read messages from array outMess
// and send it in the send of Chan for writePump func
// func (c *Client) outgoingMessagesRouter() {
//
// 	tick := time.Tick(СheckingMessages * time.Millisecond)
// 	for range tick {
// 		if !c.Status {
// 			return
// 		}
// 		if len(c.OutMess.GetMessages()) > 0 {
// 			c.Send <- c.OutMess.GetMessages()[0]
// 			c.OutMess.DelFirstM()
// 		}
// 	}
// }

// PushMessagesForRouter - push message for array letter for router
//func (c *Client) PushMessagesForRouter() {
//
//	tick := time.Tick(PushMessage * time.Millisecond)
//	for range tick {
//		if !c.Status {
//			return
//		}
//		if len(c.InMess) > 0 {
//
//			mesFor := make(Letters, 0, 500)
//			var count int
//
//			for i, val := range c.InMess {
//				lettToStr := string(val)
//				typeLett := lettToStr[0:4]
//				letter := lettToStr[4:]
//				mesFor = append(mesFor, Letter{c.ID, typeLett, letter})
//				count = i
//			}
//
//			if count == 0 {
//				(*c).InMess = (*c).InMess[1:]
//			} else {
//				(*c).InMess = (*c).InMess[count:]
//			}
//
//			LettersFrom = append(LettersFrom, mesFor...)
//
//			// (*c).InMess.DelFirstM()
//
//			// for _, val := range c.InMess {
//			// 	lettToStr := string(val)
//			// 	typeLett := lettToStr[0:4]
//			// 	letter := lettToStr[4:]
//			// 	LettersFrom.Add(Letter{c.ID, typeLett, letter})
//			// 	(*c).InMess.DelFirstM()
//			// }
//		}
//	}
//}

// start - func for start methods of client
func (c *Client) start() {

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go c.writePump()
	go c.readPump()
	// go c.outgoingMessagesRouter()
	// go c.PushMessagesForRouter()
}

// GetID - return Client ID
// type int
func (c *Client) GetID() int {
	return c.ID
}
