package verification

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func Server() {
	router := mux.NewRouter()

	router.HandleFunc("/", getTest).Methods("GET")
	router.HandleFunc("/authorization", Authorization).Methods("POST")
	router.HandleFunc("/registration", Registration).Methods("POST")
	// router.HandleFunc("/<your-url>", <function-name>).Methods("<method>")

	err := http.ListenAndServe("192.168.0.65:55442", router) // *addr
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func getTest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("TDB")
}
