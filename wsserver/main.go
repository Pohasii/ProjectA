package main

import "projecta.com/me/wsserver/src"

// define redis object
var Red src.Redis = src.Redis{}

func main() {

	// init redis connection
	Red.init()
}

