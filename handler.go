package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type requestForEncryption struct {
	ID        int
	Plaintext string
}

type requestForDecryption struct {
	ID     int
	AesKey []byte
}

/*EncryptionHandler - end point to encrypt the plaintext and return key*/
func EncryptionHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	var t requestForEncryption
	err = json.Unmarshal(body, &t)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	key, err := generateKeyAndEncrypt(t.ID, t.Plaintext)
	if err != nil {
		log.Println(err.Error())
		key = []byte(err.Error())
	}
	sendResponse(w, key)
}

/*DecryptionHandler - decrupt the stored ciphertext with the key provided*/
func DecryptionHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	var t requestForDecryption
	err = json.Unmarshal(body, &t)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()
	message, err := retrieveOriginalText(t.ID, string(t.AesKey))

	if err != nil {
		log.Print(err.Error())
		message = []byte("oops something gone wrong, we could not retrieve the content!")
	}
	sendResponse(w, message)
}

func sendResponse(w http.ResponseWriter, message []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(message)
}
