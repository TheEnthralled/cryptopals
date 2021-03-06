package main

import (
	"crypto/rand"
	"log"
	"cryptopals/utils/aes"
	"cryptopals/utils"
	"fmt"
)

const blockLen = 16

var plainTexts = [...]string {
	"MDAwMDAwTm93IHRoYXQgdGhlIHBhcnR5IGlzIGp1bXBpbmc=",
	"MDAwMDAxV2l0aCB0aGUgYmFzcyBraWNrZWQgaW4gYW5kIHRoZSBWZWdhJ3MgYXJlIHB1bXBpbic=",
	"MDAwMDAyUXVpY2sgdG8gdGhlIHBvaW50LCB0byB0aGUgcG9pbnQsIG5vIGZha2luZw==",
	"MDAwMDAzQ29va2luZyBNQydzIGxpa2UgYSBwb3VuZCBvZiBiYWNvbg==",
	"MDAwMDA0QnVybmluZyAnZW0sIGlmIHlvdSBhaW4ndCBxdWljayBhbmQgbmltYmxl",
	"MDAwMDA1SSBnbyBjcmF6eSB3aGVuIEkgaGVhciBhIGN5bWJhbA==",
	"MDAwMDA2QW5kIGEgaGlnaCBoYXQgd2l0aCBhIHNvdXBlZCB1cCB0ZW1wbw==",
	"MDAwMDA3SSdtIG9uIGEgcm9sbCwgaXQncyB0aW1lIHRvIGdvIHNvbG8=",
	"MDAwMDA4b2xsaW4nIGluIG15IGZpdmUgcG9pbnQgb2g=",
	"MDAwMDA5aXRoIG15IHJhZy10b3AgZG93biBzbyBteSBoYWlyIGNhbiBibG93",
}


var key []byte

func Encryptor() ([]byte, []byte) {
	if len(key) == 0{
		key = make([]byte, 16)
		_, err := rand.Read(key)
		if err != nil {
			log.Fatal(err)
		}
			
	}

	iv := make([]byte, 16)
	_, err := rand.Read(iv)
	if err != nil {
		log.Fatal(err)
	}
	
	i,err := utils.RandInt(0, len(plainTexts))
	if err != nil {
			log.Fatal(err)
	}
	
	
	cipherText := aes.Encrypt_CBC_PKCS7(key, iv, []byte(plainTexts[i]))	
	return cipherText, iv
}

func Oracle(cipherText []byte, iv []byte) bool {
	_, err := aes.Decrypt_CBC_PKCS7(key, iv, cipherText)
	if err != nil {
		return false
	}
	return true
}


func Attack() string {
	cipherText, iv := Encryptor()

	blocks := make([]byte, len(iv) + len(cipherText))
	copy(blocks[:len(iv)], iv)
	copy(blocks[len(iv):], cipherText)

	plainText := make([]byte, len(cipherText))
	
	for i := 0; i < len(cipherText); i += blockLen  {
		C1 :=         blocks[i : i + blockLen]
		tamperedC1 := make([]byte, blockLen); copy(tamperedC1, C1)
		C2 := blocks[i + blockLen : i + (blockLen << 1)]
		DC2 := make([]byte, blockLen)

		for p := 0; p < blockLen; p++ {
			for j := 0; j < p; j++ {
				tamperedC1[blockLen-j-1] = DC2[blockLen-j-1] ^ byte(p+1)
			}

			triedB := make([]bool, 256)
			numBtried := 0
			
			for {
				b, _ := utils.RandInt(0, 256)
				if triedB[b] {
					continue
				}
				triedB[b] = true
				numBtried++
				if numBtried == 256 {
					return ""
				}
				tamperedC1[blockLen-p-1] = byte(b)
				if Oracle(C2, tamperedC1) {
					DC2[blockLen-p-1] = byte(b) ^ byte(p+1)
					break
				}
			}	
		}	
		for j := 0; j < blockLen; j++ {
			plainText[i+j] = DC2[j] ^ C1[j]
		}	
	}
	return string(plainText)
}

// TODO: Make parallel version using context package

func main() {
	var res string
	for ; res == ""; res = Attack() {	
	}
	fmt.Println(res)
}
