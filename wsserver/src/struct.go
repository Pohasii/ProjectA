package main

// UsersOnline and UserOnline used in PushOnlineClientsToChat fucn, file connections
type UsersOnline []UserOnline

type UserOnline struct {
	ID   int    `json:"id"`
	Nick string `json:"nick"`
}
