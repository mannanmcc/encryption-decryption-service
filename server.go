package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"math/rand"
	"time"
)

type clientData struct {
	ID      int
	key     string
	payload byte
}

var iv = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05}

const keyString = "abcdefghijklmnopqrstuvwxyz1234567890"

func encrypt(key []byte, text string) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	plaintext := []byte(text)
	cfb := cipher.NewCFBEncrypter(block, iv)
	ciphertext := make([]byte, len(plaintext))
	cfb.XORKeyStream(ciphertext, plaintext)

	return encodeBase64(ciphertext), nil
}

func decrypt(key, text string) ([]byte, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}
	ciphertext := decodeBase64(text)
	cfb := cipher.NewCFBEncrypter(block, iv)
	plaintext := make([]byte, len(ciphertext))
	cfb.XORKeyStream(plaintext, ciphertext)

	return plaintext, nil
}

func encodeBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func decodeBase64(s string) []byte {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return data
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func generateKey() []byte {
	b := make([]byte, 24)
	for i := range b {
		b[i] = keyString[rand.Intn(len(keyString))]
	}

	return b
}

/*generateKeyAndEncrypt - it generate the key and encrypt the plain text*/
func generateKeyAndEncrypt(id int, plaintext string) ([]byte, error) {
	db := NewDB()
	err := db.CheckIfIDAllReadyExists(id)
	if err != nil {
		return nil, err
	}

	publicKey := generateKey()
	encryptedData, err := encrypt(publicKey, plaintext)

	if err != nil {
		return nil, err
	}
	err = db.store(id, string(publicKey), encryptedData)

	if err != nil {
		return nil, err
	}

	return publicKey, nil
}

/*RetrieveOriginalText - it retrieve content from database and return after decryption */
func retrieveOriginalText(id int, aesKey string) ([]byte, error) {
	db := NewDB()
	data, keyInDb, err := db.retrieveDecryptedContent(id)
	if keyInDb != aesKey {
		return nil, errors.New("Key provided in the request cannot be used to descrypt the cipher text")
	}
	if err != nil {
		return nil, err
	}

	return decrypt(aesKey, data)
}
