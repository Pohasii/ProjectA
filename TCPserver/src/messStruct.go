package main

type messageForUser struct {
	ConnID  int    `json:"id"`
	Message []byte `json:"message"`
}
