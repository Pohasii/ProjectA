package main

import (
	"fmt"
)

var MessageFrom chan []byte = make(chan []byte, 2000)
var MessageFor chan []byte = make(chan []byte, 2000)

// struct for save user connections
var Connections *conns = &conns{
	c: make(map[int]conn, 0),
	lastID: 0,
}

var TCPSocket *Socket = &Socket{}

func main() {
	fmt.Println("TCP Server starting..")

	SetEnv()
	go Connections.messagesRouter()

	TCPSocket.SetTcpAddr()
	TCPSocket.setTcpHandler()
	// TCPSocket.setDeadline(10000)
	go func(){

			for mess := range MessageFrom {
				fmt.Printf("Message: %v", string(mess))
				for _, val := range Connections.c {
					fmt.Println("send to: ", val.ConnID, " the mess: ", mess)
					val.Sent <- mess
				}
			}

		//for {

			// mes := messageForUser{}
			// var mess string
			// fmt.Print("Введите сообщение: ")
			// fmt.Fscan(os.Stdin, &mess)

			// fmt.Print("Введите id: ")
			// fmt.Fscan(os.Stdin, &mes.ConnID)

			// mes.Message = []byte(mess)

			// m, err := json.Marshal(mes)
			// if err != nil {
			// 	fmt.Printf("json.Marshal(mes): %v", err)
			// }
			// MessageFor <- m
		//}
	}()
	TCPSocket.acceptConn(Connections)

}
