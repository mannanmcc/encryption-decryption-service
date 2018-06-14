package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	var router = mux.NewRouter()
	router.HandleFunc("/encrypt", EncryptionHandler).Methods("POST")
	router.HandleFunc("/decrypt", DecryptionHandler).Methods("POST")

	log.Fatal(http.ListenAndServe(":8000", router))
}
