package client

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"log"
)

// Profile - user profile
type Profile struct {
	ID        int    `json:"id"`
	Nick      string `json:"nick"`
	Password  string
	Login     string
	Token     string `json:"token"`
	Friends   []int  `json:"friends"`
	IsActive  bool
	FirstConn bool
	connID int
}

type Clients map[int]Profile

func CheckToken (token string) (Profile, CheckErr) {
	DbConn.setCollection("users")

	var profile Profile
	checkErr := CheckErr{}

	err := DbConn.Collection.FindOne(DbConn.ctx, bson.D{{"token", token}}).Decode(&profile)

	if err != nil {
		log.Println(err)
		checkErr.Code = 0
		checkErr.Text = err.Error()
		return profile, checkErr
	} else {
		checkErr.Code = 1
		checkErr.Text = "Successful"
		return profile, checkErr
	}
}

func CheckNick (nick string) bool {
	DbConn.setCollection("users")

	count, err := DbConn.Collection.CountDocuments(DbConn.ctx, bson.D{{"nick", nick}})
	if err != nil {
		log.Println("err:", err)
	}

	if count > 0 {
		return false
	}

	return true

}

func SetNick () {
	DbConn.setCollection("users")

	count, err := DbConn.Collection.CountDocuments(DbConn.ctx, bson.D{{"login", Credential.Login}})
	if err != nil {
		log.Println("err:", err)
		json.NewEncoder(w).Encode(Failed{"Something went wrong. Please, try again."})
	}
}