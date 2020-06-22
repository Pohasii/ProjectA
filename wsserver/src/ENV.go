package wslib

import "os"

const (
	WsServerHost = "127.0.0.1"
	WsServerPort = "55443"
	DbServerHost = "127.0.0.1"
	DbServerPort = "27001"
	RedisHost = "localhost"
	RedisPort = "6379"
)

//
func SetEnv() {
	os.Setenv("WebsocketIP", WsServerHost)
	os.Setenv("WebsocketPORT", WsServerPort)
	os.Setenv("DataBaseIP", DbServerHost)
	os.Setenv("DataBasePORT", DbServerPort)
	os.Setenv("RedisConn", RedisHost+":"+RedisPort)
}
