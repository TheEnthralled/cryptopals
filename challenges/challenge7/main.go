package main

import (
	"crypto/aes"
	//"crypto/cipher"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func DecryptAES_ECB(key []byte, cipherText []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatal(err)
	}
	
	decryptedText := make([]byte, len(cipherText))
	blockSize := block.BlockSize()
	for i := 0; i < len(cipherText); i += blockSize {
		block.Decrypt(decryptedText[i : i + blockSize],
			      cipherText[   i : i + blockSize])
	}
	
	return decryptedText
}

func main(){
	key := []byte(os.Args[1])
	fileContents, err := ioutil.ReadFile(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}

	ciperText, err := base64.StdEncoding.DecodeString(string(fileContents))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("key: %s \n", key)
	fmt.Println("decrypted text:")
	fmt.Println(string(DecryptAES_ECB(key, ciperText)))
}
