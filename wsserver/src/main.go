package main

import (
	redis "projecta.com/me/redis"
)

// define redis object
var Red redis.Redis = redis.Redis{}

func main() {

	// init redis connection
	Red.Init(0)
	Start()
}

