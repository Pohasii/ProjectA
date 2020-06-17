package verification

import (
	"os"
)

const (
	AuthServerHost = "0.0.0.0" //"127.0.0.1"
	AuthServerPort = "80"
	DbServerHost   = "mongodb"
	DbServerPort   = "27017"
	DataBaseLogin  = "root"
	DataBasePass   = "rootpassword"
	RedisHost = "redis"
	RedisPort = "6379"
)

//
func SetEnv() {
	conn := "mongodb://"+DataBaseLogin+":"+DataBasePass+"@"+DbServerHost+":"+DbServerPort
	os.Setenv("AuthenticationIP", AuthServerHost)
	os.Setenv("AuthenticationPORT", AuthServerPort)
	os.Setenv("DBConn", conn)
	os.Setenv("RedisConn", RedisHost +":"+ RedisPort)
}
