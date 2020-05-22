package verification

import "os"

const (
	AuthServerHost = "127.0.0.1"
	AuthServerPort = "55443"
	DbServerHost   = "127.0.0.1"
	DbServerPort   = "27001"
)

//
func SetEnv() {
	os.Setenv("AuthenticationIP", AuthServerHost)
	os.Setenv("AuthenticationPORT", AuthServerPort)
	os.Setenv("DataBaseIP", DbServerHost)
	os.Setenv("DataBasePORT", DbServerPort)
}
