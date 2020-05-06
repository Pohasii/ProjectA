package wsserver

import (
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

// Clients - []Client
type Connections []Client

// Add - func for add new client for Clients array
// expected WS conn
func (c *Connections) Add(conn *websocket.Conn) {

	var Client Client = Client{
		ID:      len(*c),
		Send:    make(chan []byte, maxMessageSize),
		OutMess: make(Messages, 0, MessagesCapacity),
		InMess:  make(Messages, 0, MessagesCapacity),
	}

	Client.Conn = conn
	Client.Status = true
	*c = append((*c)[:], Client)
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
	switch id {
	case 0:
		*c = append((*c)[1:])
	case len(*c):
		*c = append((*c)[0 : id-1])
	default:
		*c = append((*c)[:id], (*c)[id+1:]...)
	}
}

// SendAll - send messages all client
func (c *Connections) SendAll(Message Message) {
	if len(*c) > 0 {
		for i := range *c {
			(*c)[i].OutMess.AddMessage(Message)
		}
	}
}

// SendOtherClient - send messages all client
func (c *Connections) SendOtherClient(ID int, Message Message) {
	if len(*c) > 0 {
		for i := range *c {
			if i != ID {
				(*c)[i].OutMess.AddMessage(Message)
			}
		}
	}
}

// SendByIDClient - send messages by ID to client
func (c *Connections) SendByIDClient(ID int, Message Message) {
	if len(*c) > 0 {
		(*c)[ID].OutMess.AddMessage(Message)
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
					// fmt.Println("before del: ", *c)
					(*c).DelByID(ID)
					// fmt.Println("after del: ", *c)
				}

				fmt.Println("remove bad connections: ", thisdel)
				thisdel = make([]int, 0, MaxConnections)
			}
		}
	}
}

// GetOfflineClient - return ids offline client
func (c *Connections) GetOfflineClient() []int {
	thisOff := make([]int, 0, MaxConnections)
	if len(*c) > 0 {
		for i := range *c {
			if (*c)[i].Status == false {
				thisOff = append(thisOff, (*c)[i].ID)
			}
		}
		return thisOff
	}
	return thisOff

}

type UsersOnline []UserOnline
type UserOnline struct {
	ID   int    `json:"id"`
	Nick string `json:"n"`
}

// GetOnlineClients - return ids offline client
func (c *Connections) GetOnlineClients() UsersOnline {
	thisOn := make(UsersOnline, 0, MaxConnections)
	if len(*c) > 0 {
		for i := range *c {
			if (*c)[i].Status == true {
				thisOn = append(thisOn, UserOnline{(*c)[i].ID, (*c)[i].Nick})
			}
		}
		return thisOn
	}
	return thisOn
}

//// CleanOffConn - remove client with off status
//func (c *Clients) CleanOffConn(status *bool, badConn []int) {
//	// tick := time.Tick(ClearPer * time.Millisecond)
//	//for range tick {
//		if len(*c) > 0 && *status {
//			thisdel := make([]int, 0, MaxConnections)
//			for i := range *c {
//				if (*c)[i].Status == false {
//					thisdel = append(thisdel, i)
//				}
//			}
//			if len(thisdel) > 0 {
//				for _, ID := range thisdel {
//					// fmt.Println("before del: ", *c)
//					(*c).DelByID(ID)
//					// fmt.Println("after del: ", *c)
//				}
//
//				fmt.Println("remove bad connections: ", thisdel)
//				thisdel = make([]int, 0, MaxConnections)
//				*status = false
//			}
//		}
//	// }
//}

//// CleanOffConn - remove client with off status
//func (c *Clients) CleanOffConn(badConn []int) {
//	// tick := time.Tick(ClearPer * time.Millisecond)
//	//for range tick {
//	if len(*c) > 0 {
//		if len(badConn) > 0 {
//			for _, ID := range badConn {
//				fmt.Println("before del: ", *c)
//				(*c).DelByID(ID)
//				fmt.Println("after del: ", *c)
//			}
//
//			fmt.Println("remove bad connections: ", badConn)
//		}
//	}
//	// }
//}
