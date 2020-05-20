package client

import (
	"fmt"
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

func (c *Profile) addToFriends (id int) {
	(*c).Friends = append((*c).Friends, id)
}

func (c *Profile) CheckFriendByID (ID int) bool {
	for _, friendID := range (*c).Friends {
		if friendID == ID {
			return true
		}
	}
	return false
}

func (c *Profile) removeFromFriends (id int) {

	friends := (*c).Friends
	fmt.Println("bef: ", friends)
	for i, val := range friends {
		if val == id {
			fmt.Println("i = ", i, " val = ", val)
			switch i {
			case 0:
				friends = append((friends)[1:])
			case len(friends)-1:
				friends = append(friends[0 : i])
			default:
				friends = append(friends[:i], friends[i+1:]...)
			}
		}
	}
	fmt.Println("aft: ", friends)
	(*c).Friends = friends
}

type Clients map[int]Profile

func (c *Clients) searchByConnID (ID int) (Profile, CheckErr) {
	result := CheckErr{}
	result.Code = 1
	result.Text = "Nothing found"
	for _, val := range *c {
		if val.connID == ID {
			return val, result
		}
	}

	result.Code = 0
	return Profile{}, result
}

func CheckToken (token string) (Profile, CheckErr) {
	DbConn.setCollection("users")

	var profile = Profile{
		Friends: make([]int,0,100),
	}
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

func SetNick (nick string, ID int) bool {
	DbConn.setCollection("users")

	request := bson.D{
		{"id", ID},
	}

	update := bson.D{
		{"$set", bson.D{
			{"nick", nick},
		}},
	}

	_, err := DbConn.Collection.UpdateOne(DbConn.ctx, request, update)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

func GetProfile (token string) GetProf {
	DbConn.setCollection("users")

	Client := GetProf{}

	res := DbConn.Collection.FindOne(DbConn.ctx, bson.D{{"token", token}}).Decode(&Client)
	if res != nil {
		log.Println("GetProfile: ", res)
	}

	return Client
}

func GetProfileByID (ID int) GetProf {
	DbConn.setCollection("users")

	Client := GetProf{}

	res := DbConn.Collection.FindOne(DbConn.ctx, bson.D{{"id", ID}}).Decode(&Client)
	if res != nil {
		log.Println("GetProfileByID: ", res)
	}
	return Client
}

func SearchByNick (nick string) (ResponseForSearchByNick, CheckErr) {
	DbConn.setCollection("users")

	var user Profile
	response := ResponseForSearchByNick{}
	checkErr := CheckErr{}

	err := DbConn.Collection.FindOne(DbConn.ctx, bson.D{{"nick", nick}}).Decode(&user)

	response.ID = user.ID
	response.Nick = user.Nick

	if err != nil {
		log.Println("SearchByNick: ", err)
		checkErr.Code = 0
		checkErr.Text = err.Error()
		return response, checkErr
	} else {
		checkErr.Code = 1
		checkErr.Text = "Successful"
		return response, checkErr
	}
}

func AddToFriends (ID int ,Friends []int) bool {
	DbConn.setCollection("users")

	request := bson.D{
		{"id", ID},
	}

	update := bson.D{
		{"$set", bson.D{
			{"friends", Friends},
		}},
	}

	_, err := DbConn.Collection.UpdateOne(DbConn.ctx, request, update)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

func removeFriends (ID int ,Friends []int) bool {
	DbConn.setCollection("users")

	request := bson.D{
		{"id", ID},
	}

	update := bson.D{
		{"$set", bson.D{
			{"friends", Friends},
		}},
	}

	_, err := DbConn.Collection.UpdateOne(DbConn.ctx, request, update)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

func removeItemFromArray (item int, array []int) []int {

	for i, val := range array {
		if val == item {
			switch i {
			case 0:
				array = append((array)[1:])
			case len(array)-1:
				array = append(array[0 : item-1])
			default:
				array = append(array[:item], array[item+1:]...)
			}
		}
	}

	return array
}