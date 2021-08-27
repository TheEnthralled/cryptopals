package main

import (
	"fmt"
	"log"
	"os"
	"encoding/hex"
	"io/ioutil"
)

func RepeatKeyXOR(key string, text string) string {
	keyLen, textLen := len(key), len(text)

	encodedText := make([]byte, textLen)
	keyIndex := 0
	for i := 0; i < textLen; i++ {
		encodedText[i] = text[i] ^ key[keyIndex]
		keyIndex = (keyIndex + 1) % keyLen
	}
	fmt.Println(string(encodedText))
	return hex.EncodeToString(encodedText)
}

func main(){
	key := os.Args[1]	
	fileBytes, err := ioutil.ReadFile(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(fileBytes))
	fmt.Println(RepeatKeyXOR(key, string(fileBytes)))
}
