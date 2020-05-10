package client

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// ChanInClient - channel for message for client service
var ChanInClient = make(chan []byte, 1000)

// ChanFromClient - channel for message from client service
var ChanFromClient = make(chan []byte, 1000)

func GetChanInClient() chan []byte {
	return ChanInClient
}

func GetChanFromClient() chan []byte {
	return ChanFromClient
}

// DbConn db conn
var DbConn mongodb

// Profiles - list clients profile
var Profiles Clients = make(Clients)

func initDB () {
	DbConn.initConnDB(os.Getenv("DataBaseIP"), "projecta", "users")
	// DbConn.initConnDB("mongodb://"+ os.Getenv("DataBaseIP")+":"+os.Getenv("DataBasePORT"), "ProjectA", "users")
}

func Start() {

	fmt.Println("Client's service started")

	initDB()
	defer DbConn.close()

	Router()
}

// Router chat logic
func Router() {
	for letter := range ChanInClient {

		let := Letter{}
		err := json.Unmarshal(letter, &let)
		if err != nil {
			log.Println(err)
		}

		switch let.LetterType {

		case "1001":
			// Check Token
			auth := Auth1001{}
			err := json.Unmarshal([]byte(let.Scroll), &auth)
			if err != nil {
				log.Println(err)
			}

			profile, checkErr := CheckToken(auth.Token)
			profile.connID = let.ClientID
			if checkErr.Code != 0 {
				Profiles[profile.ID] = profile
				send := AuthSuccessful{true}
				ChanFromClient <- toByte(Letter{let.ClientID, "1001",toJsonString(send)})
				ChanFromClient <- toByte(Letter{let.ClientID, "1902","Change status"})

				ChanFromClient <- toByte(Letter{87654321, "2550","{[]}"})
			} else {
				send := AuthSuccessful{false}
				ChanFromClient <- toByte(Letter{let.ClientID, "1001",toJsonString(send)})
			}
		case "1002":
			// Set nick
			nick := SetNick1002{}

			err := json.Unmarshal([]byte(let.Scroll), &nick)
			if err != nil {
				log.Println(err)
			}

			if CheckNick(nick.Nick) {
				for _, profile := range Profiles {
					if profile.connID == let.ClientID {
						if res := SetNick(nick.Nick,profile.ID); res {
							send := SetNick1002Succ{true}
							ChanFromClient <- toByte(Letter{let.ClientID, "1002",toJsonString(send)})
							break
						} else {
							send := SetNick1002Err{"Oops, failed"}
							ChanFromClient <- toByte(Letter{let.ClientID, "1002",toJsonString(send)})
							break
						}
					}
				}
			} else {
				send := SetNick1002Err{"This nickname is existent."}
				ChanFromClient <- toByte(Letter{let.ClientID, "1002",toJsonString(send)})
			}
		case "1003":
			//TDB
			token := Auth1001{}
			err := json.Unmarshal([]byte(let.Scroll), &token)
			if err != nil {
				log.Println(err)
			}
			prof := GetProfile(token.Token)
			ChanFromClient <- toByte(Letter{let.ClientID, "1003",toJsonString(prof)})
		case "1004":
			// Get users
		case "1901":
			//only check
			for _, profile := range Profiles {
				if profile.connID == let.ClientID {
					delete(Profiles, profile.ID)
					fmt.Println("Client, Router: remove user id:", profile.connID)
					break
				}
			}

		default:
			// ignore mess
			fmt.Println("Client, type messages incorrect: ", let.LetterType)
		}
	}
}


func toByte(letter Letter) []byte {
	send, err := json.Marshal(letter)
	if err != nil {
		log.Println(err)
	}
	return send
}

func toJsonString(letter interface{}) string {
	send, err := json.Marshal(letter)
	if err != nil {
		log.Println(err)
	}
	return string(send)
}
