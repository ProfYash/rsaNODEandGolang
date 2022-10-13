package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var Privatekey, _ = rsa.GenerateKey(rand.Reader, 2048)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/api/encrypt", generateCipherText).Methods("POST")
	r.HandleFunc("/api/decrypt", generateDecipherText).Methods("POST")
	r.HandleFunc("/api/getpublicKey", getPublicKeyFromServer).Methods("GET")
	// r.HandleFunc("/api/getpublicKey", getPublicKeyFromServer).Methods("POST")

	fmt.Println("Starting server. Running at 9002")
	log.Fatal(http.ListenAndServe(":9002", r))
}

func generateDecipherText(w http.ResponseWriter, r *http.Request) {
	fmt.Println("this sgot hit")
	var ct CipherText
	json.NewDecoder(r.Body).Decode(&ct)
	fmt.Println(ct)
	fmt.Println(" Message!!")
	fmt.Println(ct.Message)
	// hexstring := hex.EncodeToString([]byte(ct.Message))
	msg, _ := hex.DecodeString(ct.Message)
	fmt.Println(msg)
	// // fmt.Println("private key recieved :", Privatekey)
	plaintext, _ := rsa.DecryptOAEP(sha1.New(), rand.Reader, Privatekey, msg, nil)
	fmt.Println(plaintext)
	w.Write(plaintext)

}

type Key struct {
	PublicKey string `json:"n"`
	e         int    `json:"e"`
}

func getPublicKeyFromServer(w http.ResponseWriter, r *http.Request) {
	// fmt.Println(reflect.TypeOf(Privatekey.PublicKey.N))
	pubASN1, err := x509.MarshalPKIXPublicKey(&Privatekey.PublicKey)
	if err != nil {
		fmt.Println("Error here********")
	}
	pubBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubASN1,
	})
	fmt.Println(pubBytes)
	// fmt.Println(string(pubBytes))
	// json.NewEncoder(w).Encode(pubBytes)
	w.Write(pubBytes)
}

func generateCipherText(w http.ResponseWriter, r *http.Request) {
	var pt PlainText
	json.NewDecoder(r.Body).Decode(&pt)

	msg := pt.Message
	publickey := Privatekey.PublicKey
	cipher, _ := rsa.EncryptOAEP(sha1.New(), rand.Reader, &publickey, []byte(msg), nil) //encryption

	// fmt.Println("cipher:\n", string(cipher))
	cipherText := hex.EncodeToString(cipher)
	json.NewEncoder(w).Encode(cipherText)
	// w.Write([]byte(cipherText))

}

type CipherText struct {
	Message string `json:"message"`
}

type PlainText struct {
	Message string `json:"message"`
}
