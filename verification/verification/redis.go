package verification

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"log"
	"os"
	"strconv"
	"time"
)

var redisConn *redis.Client

func redisInit() {
	count := 0
	redisConn = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("RedisConn"),
		Password: "",
		DB:       0,
	})

	pong, err := redisConn.Ping().Result()
	if err != nil {
		log.Println("auth service, redisInit error: ", err)
		go func(){
			if count > 10 {
				return
			}
			count++
			time.Sleep(500)
			redisInit()
		}()
	}
	fmt.Println(pong, err)
}

func itemForRedis(body Profile) (id string, result []byte) {

	result, err := json.Marshal(body)
	if err != nil {
		log.Println("auth service, itemForRedis error: ", err)
	}

	id = "id" + strconv.Itoa(body.ID)

	return
}

func rediset(user Profile) {

	dur, err := time.ParseDuration("5m")
	if err != nil {
		log.Println(err)
	}

	key, body := itemForRedis(user)
	err = redisConn.Set(key, body, dur).Err()
	if err != nil {
		log.Println("auth service, rediset error: ", err)
	}

}
