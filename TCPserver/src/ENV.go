package main

import "os"

const (
	TCPServerHost = "127.0.0.1"
	TCPServerPort = "55443"
	RedisHost     = "localhost"
	RedisPort     = "6379"
)

func SetEnv() {
	os.Setenv("TCPConn", TCPServerHost+":"+TCPServerPort)
	os.Setenv("RedisConn", RedisHost+":"+RedisPort)
}
