package redis

import (
	"fmt"
	"github.com/go-redis/redis"
	"log"
	"os"
	"time"
)

// new version

// Redis - object redis
type Redis struct {
	conn *redis.Client
}

// init connection with Redis server
func (r *Redis) Init(numbDB int) {
	count := 0
	r.conn = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("RedisAddr"),
		Password: os.Getenv("RedisPass"),
		DB:       numbDB,
	})

	pong, err := r.conn.Ping().Result()
	if err != nil {
		log.Println("wsServer service, redisInit error: ", err)
		go func(){
			if count > 10 {
				return
			}
			count++
			time.Sleep(1000)
			r.Init(numbDB)
		}()
	}
	fmt.Println(pong, err)
}

func (r *Redis) PushToRedis(key string, data []byte) (error, bool) {

	dur, err := time.ParseDuration("60m")
	if err != nil {
		log.Println(err)
	}

	err = r.conn.Set(key, data, dur).Err()
	if err != nil {
		log.Println("wsServer service, pushToRedis error: ", err)
		return err, true
	}

	return nil, true
}

func (r *Redis) Get (id string) (res []byte){

	res, err := r.conn.Get(id).Bytes() //Result()
	if err != nil {
		log.Println("wsServer service, (r *Redis) get->Get error: ", err)
	}

	return
}