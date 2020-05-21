package client

import "os"

const (
	DbServerHost = "127.0.0.1"
	DbServerPort = "27001"
)

//
func SetEnv() {
	os.Setenv("DataBaseIP", DbServerHost)
	os.Setenv("DataBasePORT", DbServerPort)
}