package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
	//"net/url"
	"crypto/rand"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

const KEY_LEN = 16
var KEY []byte


func genKey() {
	KEY = make([]byte, KEY_LEN)
	_, err := rand.Read(KEY)
	if err != nil {
		log.Fatal(err)
	}
}


func validateSignature(signature []byte, message []byte) bool {
	mac := hmac.New(sha256.New, KEY)
	mac.Write(message)
	sig2 :=  hex.EncodeToString(mac.Sum(nil))
	
	if len(signature) == 0 || len(signature) != len(sig2) {
		return false
	}
	
	for i := 0; i < len(signature); i++ {
		if sig2[i] != signature[i] {
			return false
		}
		time.Sleep(50 * time.Millisecond)
	}
	
	return true
}

func handler(w http.ResponseWriter, r *http.Request) {	
	q := r.URL.Query()
	file := q.Get("file")
	signature := q.Get("signature")

	//fmt.Println("signature = ", signature)
	
	if len(KEY) == 0 {
		genKey()
	}

	res := validateSignature([]byte(signature), []byte(file))
	if res {
		w.WriteHeader(200)
		fmt.Fprintf(w, "ACCESS GRANTED")
	} else {
		w.WriteHeader(500)
		fmt.Fprintf(w, "ACCESS DENIED: Invalid MAC")
	}
}

func Server() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":9000", nil))
}
