package verification

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
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

		connDb.initConnDB(os.Getenv("DBConn"), "ProjectA", "users")
		defer connDb.close()

		if checkAkk(connDb, Credential) {
			client, status := updateAkk(connDb, Credential)
			if status {
				json.NewEncoder(w).Encode(Successful{client.Token})
				rediset(client)
			} else {
				json.NewEncoder(w).Encode(Failed{"oops, failed to complete authorization"})
			}
		} else {
			json.NewEncoder(w).Encode(Failed{"Sorry, your password or login is incorrect"})
		}
	}
}
