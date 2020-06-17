package wsserver

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
		log.Println("wsServer service, redisInit error: ", err)
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

type Profile struct {
	ID        int    `json:"id"`
	Nick      string `json:"nick"`
	Password  string
	Login     string
	Token     string `json:"token"`
	Friends   []int  `json:"friends"`
	IsActive  bool
	FirstConn bool
	connID int
}


func itemForRedis(body Profile) (id string, result []byte) {

	result, err := json.Marshal(body)
	if err != nil {
		log.Println("wsServer service, itemForRedis error: ", err)
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
		log.Println("wsServer service, rediset error: ", err)
	}
}

func rediget(id string) (user Profile){

	res, err := redisConn.Get(id).Result()
	if err != nil {
		log.Println("wsServer service, rediget->Get error: ", err)
	}

	err = json.Unmarshal([]byte(res), &user)
	if err != nil {
		log.Println("wsServer service, rediget->Unmarshal error: ", err)
	}

	return
}