package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// new

// Client Hub
type ClientHub struct {
	Connections map[int]Client
	ID          ID
}

// getID () int - return last client ID
// ch.ID++
// return ch.ID
func (ch *ClientHub) getID() int {
	return ch.ID.Get()
}

func (ch *ClientHub) AddClient(client Client) bool {
	ID := ch.getID()

	// _, ok := m["route"]
	_, ok := ch.Connections[ID]
	if ok {
		ch.Connections[ID] = client
		return true
	}

	return false
}

// DelClientByID - DelClientByID(ID int) bool
func (ch *ClientHub) DelClientByID(ID int) bool {

	delete(ch.Connections, ID)
	_, ok := ch.Connections[ID]
	if ok == false {
		return true
	}
	return false
}

// GetClientByID(ID int) (bool, Client)
func (ch *ClientHub) GetClientByID(ID int) (bool, Client) {

	client, ok := ch.Connections[ID]
	return ok, client
}

// struct ID ===============================================================<
// 	sync.RWMutex
//	ID int
type ID struct {
	sync.RWMutex
	ID int
}

// Get () int
// safe with Mutex
func (id *ID) Get() int {
	// block for read
	id.RLock()

	// change ID
	id.ID++

	// unblock
	id.RUnlock()

	// return result
	return id.ID
}

// ========= end struct ID =====================================>






// OLD
// OLD
// OLD
// OLD
// OLD
// OLD
// OLD
// OLD
// OLD
// OLD
// OLD
// OLD
// OLD
// OLD ========================================================

// Connections - []Client
type Connections []Client

// id for user
var lastID int

// Add - func for add new client for Clients array
// expected WS conn
func (c *Connections) Add(conn *websocket.Conn) {

	var Client Client = Client{
		ID:   lastID,
		Send: make(chan []byte, maxMessageSize),
	}

	Client.Conn = conn
	Client.Status = true
	*c = append(*c, Client)
	lastID++
	go func() {
		time.Sleep(TimeForAuth * time.Millisecond)
		if (*c)[Client.ID].Auth == false {
			(*c)[Client.ID].Conn.Close()
			(*c)[Client.ID].Status = false
			if (*c)[Client.ID].Send != nil {
				close((*c)[Client.ID].Send)
			}
		}
	}()
}

// GetClients - func GetClients() Clients
// return array of *Clients
func (c *Connections) GetClients() Connections {
	return *c
}

// GetClientsID - func GetClients()
// return []int
func (c *Connections) GetClientsID(ID int) []int {

	client := make([]int, 0, len(*c))

	if len(*c) > 0 {
		for _, cl := range *c {
			if cl.Status == true {
				client = append(client, cl.ID)
			}
		}
	}
	return client // strings.Join(client, " ")
}

// DelByID - delete client from Clients array by id
// expected id (int)
func (c *Connections) DelByID(id int) {
	for i, val := range *c {
		if val.ID == id {
			switch i {
			case 0:
				*c = append((*c)[1:])
			case len(*c) - 1:
				*c = append((*c)[0:i])
			default:
				*c = append((*c)[:i], (*c)[i+1:]...)
			}
		}
	}

}

// CleanOffConn - remove client with off status
func (c *Connections) CleanOffConn() {

	tick := time.Tick(ClearPer * time.Millisecond)

	for range tick {
		if len(*c) > 0 {
			thisdel := make([]int, 0, MaxConnections)
			for i := range *c {
				if (*c)[i].Status == false {
					thisdel = append(thisdel, i)
				}
			}
			if len(thisdel) > 0 {
				for _, ID := range thisdel {
					(*c).DelByID(ID)
				}

				fmt.Println("remove bad connections: ", thisdel)
				thisdel = make([]int, 0, MaxConnections)
			}
		}
	}
}

// GetOfflineClient - return ids offline client
func (c *Connections) GetOfflineClient() []int {
	TheseOff := make([]int, 0, MaxConnections)
	if len(*c) > 0 {
		for i := range *c {
			if (*c)[i].Status == false {
				TheseOff = append(TheseOff, (*c)[i].ID)
			}
		}
		return TheseOff
	}
	return TheseOff

}
