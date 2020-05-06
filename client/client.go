package client

//import (
//	"fmt"
//	ws "gameserver/wsserver"
//	"time"
//)
//
//// Client - as
//type Client struct {
//	ID     int
//	Status bool
//	Nick   string
//}
//
//// Clients - []Client
//type Clients []Client
//
//// Add add client in system
//func (c *Clients) Add(ID int, nick string) bool {
//
//	var Client Client = Client{
//		ID:     ID,
//		Status: true,
//		Nick:   nick,
//	}
//
//	*c = append((*c)[:], Client)
//
//	return true
//}
//
//// DelByID - delete client from Clients array by id
//// expected id (int)
//func (c *Clients) DelByID(id int) {
//	switch id {
//	case 0:
//		*c = append((*c)[1:])
//	case len(*c):
//		*c = append((*c)[0 : id-1])
//	default:
//		*c = append((*c)[:id], (*c)[id+1:]...)
//	}
//}
//
//// CleanOffConn - remove client with off status
//func (c *Clients) CleanOffConn(remove []int) bool {
//
//	if len(*c) > 0 && len(remove) > 0 {
//
//		for i := range *c {
//
//			for _, val := range remove {
//
//				if (*c)[i].ID == val {
//
//					(*c).DelByID(val)
//
//				}
//			}
//		}
//		return true
//	}
//	return false
//}
//
//func (c *Clients) Cleaning(conn *ws.Clients) {
//	ticker := time.Tick(ClearPer * time.Millisecond)
//	for range ticker {
//		fmt.Println("test rtr")
//		badconn := conn.GetOfflineClient()
//		fmt.Println("badconn: ", badconn)
//		fmt.Println("Clients: ", *c)
//		clean := c.CleanOffConn(badconn)
//		fmt.Println("clean: ", clean)
//		if clean {
//			ws.RunCleanConn = clean
//			ws.Connections.CleanOffConn(badconn)
//			fmt.Println("clean: ", clean)
//		}
//	}
//}
//
