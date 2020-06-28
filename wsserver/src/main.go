package main

import (
	"fmt"
	redis "projecta.com/me/redis"
)

// define redis object
var Red redis.Redis = redis.Redis{}

var Hub ClientHub = ClientHub{
	Connections: make(map[int]*Client, 0),
}

func main() {

	// init redis connection
	// Red.Init(0)
	go func(){
		for key := range FromConnChan {
			fmt.Println(string(key))
		}
	}()
	Start()
}

