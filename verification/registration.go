package verification

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/bson"
)

func Registration(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var connDb mongodb
	var Credential Credentials

	err := json.NewDecoder(r.Body).Decode(&Credential)
	if err != nil {
		log.Println("Registration: ", err)
	}

	if Credential.Login == "" && Credential.Password == "" {
		json.NewEncoder(w).Encode(Failed{"Sorry, your password or login incorrect"})
	} else {

		connDb.initConnDB("mongodb://"+ os.Getenv("DataBaseIP")+":"+os.Getenv("DataBasePORT"), "ProjectA", "users")
		defer connDb.close()

		count, err := connDb.Collection.CountDocuments(connDb.ctx, bson.D{{"login", Credential.Login}})
		if err != nil {
			log.Println("err:", err)
			json.NewEncoder(w).Encode(Failed{"Something went wrong. Please, try again."})
		}

		if count > 0 {
			json.NewEncoder(w).Encode(Failed{"Sorry, this login is used."})
		} else {
			ids := LastID{}

			update := bson.D{
				{"$inc", bson.D{
					{"value", 1},
				}},
			}

			connDb.setCollection("userid")
			connDb.Collection.FindOneAndUpdate(connDb.ctx, bson.D{{"id", "id"}}, update).Decode(&ids)
			connDb.setCollection("users")

			newClient := Profile{
				ID:       ids.Value,
				Login:    Credential.Login,
				Password: Credential.Password,
				IsActive: true,
				Friends:  make([]int, 0),
			}

			_, err := connDb.Collection.InsertOne(connDb.ctx, newClient)
			if err != nil {
				log.Println(err)
				json.NewEncoder(w).Encode(Failed{"oops, failed to complete registration"})
			} else {
				json.NewEncoder(w).Encode(Successful{"Congratulations, you have successfully registered!"})
			}
		}
	}
}
