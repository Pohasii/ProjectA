package verification

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

func Server() {
	router := mux.NewRouter()

	router.HandleFunc("/", mainPage).Methods("GET")
	router.HandleFunc("/authorization", Authorization).Methods("POST")
	router.HandleFunc("/registration", Registration).Methods("POST")
	// router.HandleFunc("/<your-url>", <function-name>).Methods("<method>")

	fmt.Println("Verification's service started at the ", os.Getenv("AuthenticationIP")+":"+os.Getenv("AuthenticationPORT"))

	err := http.ListenAndServe(os.Getenv("AuthenticationIP")+":"+os.Getenv("AuthenticationPORT"), router) // *addr
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("TDB")
}
