package main

import (
	"cryptopals/utils/aes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func main(){
	key := "YELLOW SUBMARINE"
	iv := make([]byte, 16)
	
	fileContents , err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	cipherText, err := base64.StdEncoding.DecodeString(string(fileContents))
	if err != nil {
		log.Fatal(err)
	}
	plainText, err := aes.Decrypt_CBC_PKCS7([]byte(key), iv, cipherText)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(plainText))
}
