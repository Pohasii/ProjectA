package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

const (
	// bytes
	WriteBuffer = 5000
	ReadBuffer  = 5000

	MessagesBuffer = 500

	// terminationСompletedСonnections = 120 * time.Second
	// Millisecond
	ReadDeadline = 2000
	// WriteDeadline = 2000

	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

type conns struct {
	c      map[int]conn
	lastID int
}

func (c *conns) delByID(ID int) {

	_, ok := (*c).c[ID]
	if ok {
		delete((*c).c, ID)
	}
}

// return True if ID is free
// and return false when ID is exist
func (c *conns) checkKey(id int) bool {
	_, ok := (*c).c[id]
	return !ok
}

func (c *conns) addConn(tcpConn *net.TCPConn) error {
	// check that ID is available
	newConn := &conn{}
	if c.checkKey(c.lastID + 1) {

		// add conn
		(*c).c[c.lastID+1] = newConn.init(c.lastID+1, tcpConn, c)

		log.Printf("Add new conn adr:%v, id %v", tcpConn.RemoteAddr(), (*c).c[c.lastID+1].ConnID)
		// up ID
		(*c).lastID = c.lastID + 1
		return nil
	}
	return errors.New("TCPserver, addConn: user add failed / ID didn't available")
}

func (c *conns) messagesRouter() {
	for mess := range MessageFor {

		m := messageForUser{}
		err := json.Unmarshal(mess, &m)
		if err != nil {
			log.Printf("TCPserver, messagesRouter: %v", err)
		} else {
			_, ok := c.c[m.ConnID]
			if ok {
				fmt.Println("messagesRouter : here")
				c.c[m.ConnID].Sent <- m.Message
			}
		}
	}
}

// func (c *conns) ClearnOffConn() {
// 	ticker := time.NewTicker(terminationСompletedСonnections)
//
// 	for range ticker.C {
// 			removeConn := make([]int, 0)
// 			for _, conn := range c.c {
// 				if !conn.ActiveStatus {
// 					delete(c.c, conn.ConnID)
// 					removeConn = append(removeConn, conn.ConnID)
// 				}
// 			}
//
// 			log.Println("remove connections - IDs: ", removeConn)
// 	}
// }

//  =================================================
//  =================================================
//  =================================================

type conn struct {
	Conn         *net.TCPConn
	Sent         chan []byte
	ConnID       int
	ActiveStatus bool
	linkOnArray  *conns
}

func (c *conn) init(ID int, Conn *net.TCPConn, cs *conns) conn {
	c.Conn = Conn
	c.Sent = make(chan []byte, MessagesBuffer)
	c.ActiveStatus = true
	c.ConnID = ID
	c.linkOnArray = cs

	go c.Start()

	return *c
}

func (c *conn) uninit() {
	if err := c.Conn.Close(); err != nil {
		log.Printf("TCPserver, uninit: %v", err)
	}
	close(c.Sent)

	log.Printf("remove offline connection, ID: %v", c.ConnID)
	c.linkOnArray.delByID(c.ConnID)
}

func (c *conn) setWriteConnParam() {

	err := c.Conn.SetWriteBuffer(WriteBuffer)
	if err != nil {
		log.Printf("TCPserver, setWriteConnParam, SetWriteBuffer: %v", err)
	}

	err = c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
	if err != nil {
		log.Printf("TCPserver, setWriteConnParam, SetWriteDeadline: %v", err)
	}
}

func (c *conn) setReadConnParam() {

	err := c.Conn.SetReadBuffer(ReadBuffer)
	if err != nil {
		log.Printf("TCPserver, setReadConnParam, SetReadBuffer: %v", err)
	}

	err = c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		log.Printf("TCPserver, setReadConnParam, SetReadDeadline: %v", err)
	}
}

func (c *conn) listener(status chan int) {

	c.Conn.SetKeepAlive(true)
	c.Conn.SetKeepAlivePeriod()
	c.Conn.SetLinger()
	c.Conn.SyscallConn()

	status <- 1
	defer func() {
		// c.Conn.Close()
		// close(c.Sent)
		_, ok := <- status
		if ok {
			c.ActiveStatus = false
			c.uninit()
			close(status)
		}
	}()

	c.setReadConnParam()

	for {

		bufferBytes, err := bufio.NewReader(c.Conn).ReadBytes('\n')

		if err != nil {
			log.Printf("tyt TCPserver, listener: %v", err)
			one := make([]byte, 1)
			if  _, err := c.Conn.Read(one); err == io.EOF {
				log.Printf(" and Tyt TCPserver, listener: %v", err)
				break
			}

			break
		} else {
			MessageFrom <- bufferBytes
		}
	}
	
	fmt.Println("listener exit")
}

func (c *conn) sender(status chan int) {
	status <- 1
	
	defer func() {
		_, ok := <- status
		if ok {
			c.ActiveStatus = false
			c.uninit()
			close(status)
		}
	}()

	ticker := time.NewTicker(pingPeriod)
	for {
		select {
		case mess, ok := <- (*c).Sent:

			fmt.Println("send to client id:",c.ConnID, "Message: ", string(mess))

			_, connStatus := <- status
			if !connStatus {
				fmt.Println("_, connStatus := <- status - works 223 line")
				return
			}

			if !ok {
				// The hub closed the channel.
				_, err := c.Conn.Write([]byte("0"))
				if err != nil {
					log.Printf("TCPserver, sender: %v", err)
				}
				return
			}

			c.setWriteConnParam()

			if len(mess) > WriteBuffer {
				severalPartOfMessForSend := cutMessages(mess)
				// sender := bufio.NewWriter(c.Conn)
				for i, part := range severalPartOfMessForSend {
					byteSent, err := c.Conn.Write(part)
					if err != nil {
						log.Printf("TCPserver, sender: %v part didn't send: %v and bytes %v", err, i, byteSent)
					}
				}
			} else {
				byteSent, err := c.Conn.Write(mess)
				if err != nil {
					log.Printf("TCPserver, sender: %v didn't send bytes %v", err, byteSent)
				}
			}

		case <-ticker.C:
			c.setWriteConnParam()
			if _, err := c.Conn.Write([]byte("1")); err != nil {
				log.Printf("TCPserver, sender: %v", err)
				return
			}
		}
	}

	fmt.Println("sender exit")
}

func (c *conn) Start() {

	fmt.Println("Start here")
	status := make(chan int, 3)
	go c.listener(status)
	go c.sender(status)
	(*c).ActiveStatus = true
}

// type inits interface {
// 	Start()
// 	sender()
// 	listener()
// 	setReadConnParam()
// 	setWriteConnParam()
// 	uninit()
// }

// cutMessages for send to client
func cutMessages(mess []byte) [][]byte {

	newMes := make([][]byte, 0)

	if len(mess) > WriteBuffer {
		newMes = append(newMes, mess[:WriteBuffer])
		if len(mess) > WriteBuffer {
			newMes = append(cutMessages(mess[WriteBuffer:]), mess[:WriteBuffer])
		}
	}

	return newMes
}
