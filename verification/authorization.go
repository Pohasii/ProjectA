package verification

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func Authorization(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var connDb mongodb
	var Credential Credentials

	err := json.NewDecoder(r.Body).Decode(&Credential)
	if err != nil {
		log.Println("Authorization: ", err)
	}

	if Credential.Login == "" && Credential.Password == "" {
		json.NewEncoder(w).Encode(Failed{"Sorry, your password or login is incorrect"})
	} else {

		connDb.initConnDB("mongodb://localhost:27017", "ProjectA", "users")
		defer connDb.close()

		request := bson.D{
			{"password", Credential.Password},
			{"login", Credential.Login},
		}

		Client := Profile{}

		res := connDb.Collection.FindOne(connDb.ctx, request).Decode(&Client)
		if res != nil {
			log.Println("res:", res)
			json.NewEncoder(w).Encode(Failed{"Sorry, your password or login is incorrect"})
		} else {

			h := md5.New()
			time := strconv.Itoa(int(time.Now().Unix()))
			Client.Token = hex.EncodeToString(h.Sum([]byte(time)))

			update := bson.D{
				{"$set", bson.D{
					{"token", Client.Token},
				}},
			}

			_, err := connDb.Collection.UpdateOne(connDb.ctx, request, update)
			if err != nil {
				log.Println(err)
				json.NewEncoder(w).Encode(Failed{"oops, failed to complete authorization"})
			} else {
				json.NewEncoder(w).Encode(Successful{Client.Token})
			}
		}
	}

}
