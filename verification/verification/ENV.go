package verification

import "os"

const (
	AuthServerHost = "0.0.0.0" //"127.0.0.1"
	AuthServerPort = "80"
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
