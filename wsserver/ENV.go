package wsserver

import "os"

const (
	WsServerHost = "127.0.0.1"
	WsServerPort = "55443"
	DbServerHost = "127.0.0.1"
	DbServerPort = "27001"
)

//
func SetEnv() {
	os.Setenv("WebsocketIP", WsServerHost)
	os.Setenv("WebsocketPORT", WsServerPort)
	os.Setenv("DataBaseIP", DbServerHost)
	os.Setenv("DataBasePORT", DbServerPort)
}
