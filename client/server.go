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
	//"mongodb://localhost:27017/?readPreference=primary&appname=MongoDB%20Compass&ssl=false"
	// DbConn.initConnDB( "mongodb://localhost:27017/?readPreference=primary&appname=MongoDB%20Compass&ssl=false", "projecta", "users")
	DbConn.initConnDB("mongodb://"+ os.Getenv("DataBaseIP")+":"+os.Getenv("DataBasePORT"), "ProjectA", "users")
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
				if profile.Nick != "" {
					onl := Online{let.ClientID, profile.ID, profile.Nick}
					ChanFromClient <- toByte(Letter{87654321, "2550",toJsonString(onl)})
				}
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

							onl := Online{let.ClientID, profile.ID, profile.Nick}
							ChanFromClient <- toByte(Letter{87654321, "2550",toJsonString(onl)})
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
			// load profile
			token := Auth1001{}
			err := json.Unmarshal([]byte(let.Scroll), &token)
			if err != nil {
				log.Println(err)
			}
			prof := GetProfile(token.Token)
			ChanFromClient <- toByte(Letter{let.ClientID, "1003",toJsonString(prof)})
		case "1004":
			// Get users by Nick
			req := GetByNick{}
			err := json.Unmarshal([]byte(let.Scroll), &req)
			if err != nil {
				log.Println("Client, 1004: ", err)
			}

			user, checkErr := SearchByNick(req.Nick)
			if checkErr.Code != 0 {

				if CurUser, checkErr := Profiles.searchByConnID(let.ClientID); (CurUser.Nick == user.Nick) && checkErr.Code != 0 {
					log.Println("the user search himself.")
					continue
				}

				ChanFromClient <- toByte(Letter{let.ClientID, "1004",toJsonString(user)})
			} else {
				ChanFromClient <- toByte(Letter{let.ClientID, "1004",toJsonString(ErrorPattern{"Nothing found"})})
			}
		case "1005":
			// add to friends

			//fmt.Println("Profiles: ", Profiles)

			req := FromFriendRequest{}
			err := json.Unmarshal([]byte(let.Scroll), &req)
			if err != nil {
				log.Println("Client, 1005: ", err)
			}

			user, checkErr := Profiles.searchByConnID(let.ClientID)

			if (user.ID == req.ID) && checkErr.Code != 0 {
				log.Println("user id for add to friend == id the user.")
				continue
			}

			if user.CheckFriendByID(req.ID) {
				log.Println("the friend added")
				continue
			}

			// если true, то это запрос на  добавление
			if req.Request {
				// если это запрос на добавление, то мы находим того, кого хотят добавить и отправляем запрос
				user, checkErr := Profiles.searchByConnID(let.ClientID)
				if checkErr.Code != 0 {
					send := RequestFriend{user.ID, user.Nick, true}
					ChanFromClient <- toByte(Letter{Profiles[req.ID].connID, "1005",toJsonString(send)})
				}
			} else {
				//  если  это  не запрос - значит ответ :) мы выполняем добавление обоим участникам соглашения в друзья
				user, checkErr := Profiles.searchByConnID(let.ClientID)

				if checkErr.Code != 0 {

					status := struct{
						requester bool
						responder bool
					}{}
					// добавляем в друзья тому кто запрашивал
					if _, ok := Profiles[req.ID]; ok{

						prof := Profiles[req.ID]
						prof.addToFriends(user.ID)
						Profiles[req.ID] = prof
						status.requester = AddToFriends(req.ID, Profiles[req.ID].Friends)
					}

					// тому кто отвечал на запрос
					if _, ok := Profiles[user.ID]; ok{

						prof := Profiles[user.ID]
						prof.addToFriends(req.ID)
						Profiles[user.ID] = prof
						// Profiles[user.ID].Friends = append(Profiles[user.ID].Friends, req.ID)
						status.responder = AddToFriends(user.ID, Profiles[user.ID].Friends)
					}

					if status.requester && status.responder {
						send := RequestFriend{user.ID, user.Nick, false}
						ChanFromClient <- toByte(Letter{Profiles[req.ID].connID, "1005",toJsonString(send)})

						send = RequestFriend{req.ID, Profiles[req.ID].Nick, false}
						ChanFromClient <- toByte(Letter{Profiles[user.ID].connID, "1005",toJsonString(send)})
					}
				}
			}
		case "1006":
			// remove from friends

			// расшифровка сообщения от юзера
			// для удаление этого юзера из списка друзей.
			requestForDel := removeFriend{}

			err := json.Unmarshal([]byte(let.Scroll), &requestForDel)
			if err != nil {
				log.Println("Client, 1006: ", err)
			}

			// получить профиль пользователя который сделал запрос на удаление

			//профиль из которого  нужно удалить
			user, checkErr := Profiles.searchByConnID(let.ClientID)
			if checkErr.Code != 0 {
				// если он онлайн

				//удаляем из массива засранца
				user.removeFromFriends(requestForDel.ID)
				//обноввляем текущий объект пользователя
				Profiles[user.ID] = user
				//удаляем из базы
				removeFriends(user.ID, user.Friends)

			}

			//отправляем результат тому, кто запрос  делал
			send := friends{Profiles[user.ID].Friends}
			ChanFromClient <- toByte(Letter{let.ClientID, "1007",toJsonString(send)})

// =============================================================================
			//теперь нужно удалить у  того,  кого удалили
			//получаем его профиль из базы
			removedProfile := GetProfileByID(requestForDel.ID)

			// удаляем того, от кого приходил запрос на удаления у того, кого удалили в первой  серии:))
			removedProfile.Friends = removeItemFromArray(user.ID, removedProfile.Friends)
			// удаляем в базе биомусор неуважительный
			removeFriends(removedProfile.ID,removedProfile.Friends)

			if _, ok := Profiles[requestForDel.ID]; ok{
				// если он онлайн
				// удалил из массива онлайн
				profile := Profiles[requestForDel.ID]
				profile.removeFromFriends(user.ID)
				Profiles[requestForDel.ID] = profile

				send = friends{removedProfile.Friends}
				ChanFromClient <- toByte(Letter{profile.connID, "1007",toJsonString(send)})
			}
		case "1007":
			//  список  друзей
			user, checkErr := Profiles.searchByConnID(let.ClientID)
			if checkErr.Code != 0 {
				send := friends{user.Friends}
				fmt.Println("1007 send: ", send)
				ChanFromClient <- toByte(Letter{let.ClientID, "1007",toJsonString(send)})
			}
		case "1901":
			//only check
			for _, profile := range Profiles {
				if profile.connID == let.ClientID {
					delete(Profiles, profile.ID)

					fmt.Println("Client, Router: remove user id:", profile.ID)

					onl := Online{profile.connID, profile.ID, profile.Nick}
					ChanFromClient <- toByte(Letter{87654321, "2550",toJsonString(onl)})
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
