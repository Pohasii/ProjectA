package client

import (
	"encoding/json"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

// ChanInClient - channel for message for client service
var ChanInClient chan []byte = make(chan []byte, 1000)

// ChanFromClient - channel for message from client service
var ChanFromClient chan []byte = make(chan []byte, 1000)

// DbConn db conn
var DbConn mongodb

// Profiles - list clients profile
var Profiles Clients

func initDB () {
	DbConn.initConnDB("mongodb://localhost:27017", "ProjectA", "users")
}

func Start() {

	initDB()
	defer DbConn.close()

	var usera Profile

	err := DbConn.Collection.FindOne(DbConn.ctx, bson.D{{"login", "login"}}).Decode(&usera)
	if err != nil {
		log.Println(err)
	}

	// fmt.Println(usera.Token)
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
			ToStruct(let.Scroll, &auth)

			profile, checkErr := CheckToken(auth.Token)
			profile.connID = let.ClientID
			if checkErr.Code != 0 {
				Profiles[profile.ID] = profile
				send := AuthSuccessful{true}
				ChanFromClient <- toByte(Letter{let.ClientID, "1001",toJsonString(send)})
				ChanFromClient <- toByte(Letter{let.ClientID, "1902","Change status"})
			} else {
				send := AuthSuccessful{false}
				ChanFromClient <- toByte(Letter{let.ClientID, "1001",toJsonString(send)})
			}

		case "1002":
			// Set nick
			nick := SetNick1002{}
			ToStruct(let.Scroll,&nick)

		case "1003":
			//TDB
	

		case "1550":
			//only check


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

func ToStruct(todo, value) {
	err := json.Unmarshal([]byte(todo), value)
	if err != nil {
		log.Println(err)
	}
}
