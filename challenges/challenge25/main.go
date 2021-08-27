package main

import (
	"crypto/rand"
	"cryptopals/utils/aes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

const (
	BLOCK_SIZE = 16
	NONCE = 0
	COUNTER_START = 0
)

var KEY []byte


func GenKey() {
	KEY = make([]byte, BLOCK_SIZE)
	_, err := rand.Read(KEY)
	if err != nil {
		log.Fatal(err)
	}
}

func EnryptPlaintext() []byte {
	GenKey()
	plainText, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	
	return aes.Encrypt_CTR(0, 0, KEY, plainText)
}

func edit(cipherText []byte, key []byte, offset int, newText []byte) {
	start := (offset/BLOCK_SIZE)*BLOCK_SIZE
	end   := ((offset + len(newText) - 1)/BLOCK_SIZE + 1) * BLOCK_SIZE

	counter := uint64(start/BLOCK_SIZE)
	plainText := aes.Decrypt_CTR(NONCE, counter, key, cipherText[start:end])
	copy(plainText[offset - start : offset - start + len(newText)], newText)
	ciTxt := aes.Encrypt_CTR(NONCE, counter, key, plainText)
	copy(cipherText[start:end], ciTxt)
}

func Edit(cipherText []byte, offset int, newText []byte) {
	edit(cipherText, KEY, offset, newText)
}

func DecryptCipherText(cipherText []byte) []byte{
	dummyBlock := make([]byte, BLOCK_SIZE)
	prevCipherTextBlock := make([]byte, BLOCK_SIZE)
	newCipherTextLen := (len(cipherText)/BLOCK_SIZE + 1)*BLOCK_SIZE 
	oldCipherTextLen := len(cipherText)
	for i := 0; i < newCipherTextLen - oldCipherTextLen; i++ {
		cipherText = append(cipherText, 0)
	}
	plainText := make([]byte, newCipherTextLen)
	
	for i := 0; i + BLOCK_SIZE <= len(cipherText); i += BLOCK_SIZE {	
		copy(prevCipherTextBlock, cipherText[i : i + BLOCK_SIZE])
		Edit(cipherText, i, dummyBlock)
		newCipherTextBlock := cipherText[i : i + BLOCK_SIZE]

		for j := 0; j < BLOCK_SIZE; j++ {
			plainText[i + j] = prevCipherTextBlock[j] ^ newCipherTextBlock[j]
		}
		
	}
	return plainText[0 : oldCipherTextLen]
}

func main() {
	cipherText := EnryptPlaintext()
	fmt.Println(string(DecryptCipherText(cipherText)))
}
