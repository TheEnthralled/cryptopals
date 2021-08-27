package main

import (
	"crypto/rand"
	"cryptopals/utils/aes"
	"errors"
	"fmt"
	"log"	
)

const BLOCKLEN = 16

var KEY []byte


var CIPHERTEXT []byte

func ASCII_Compliant(text string) bool {

	for _, c := range text {
		if int(c) < 32 {
			return false
		}
	}

	return true
}

func EncryptText(text string) []byte {
	if len(KEY) == 0 {
		KEY = make([]byte, BLOCKLEN)
		_, err := rand.Read(KEY)
		if err != nil {
			log.Fatal(err)
		}	
	}

	return aes.Encrypt_CBC_PKCS7(KEY, KEY, []byte(text))
}

func DecryptText(ciphertext []byte) ([]byte, error){
	plaintext, err := aes.Decrypt_CBC_PKCS7(KEY, KEY, ciphertext)
	if err != nil && err.Error() != "pkcs: some pad byte does not equal to the total number of padded bytes"{
		return plaintext, err
	}

	if !ASCII_Compliant(string(plaintext)) {
		return plaintext, errors.New("non-ascii compliant")
	}

	return plaintext, nil
}

func AttackerFindKey(ciphertext []byte) ([]byte, error) {
	ct := make([]byte, 4*BLOCKLEN)
	copy(ct[:BLOCKLEN], ciphertext[:BLOCKLEN])
	copy(ct[2*BLOCKLEN:3*BLOCKLEN], ciphertext[:BLOCKLEN])
	copy(ct[3*BLOCKLEN:], ciphertext[3*BLOCKLEN:])
	
	plaintext, err := DecryptText(ct)
	
	if err != nil && err.Error() != "non-ascii compliant" {
			return nil, err
	} else if err == nil {
		return nil, errors.New("couldn't get desied plantext")
	}

	key := make([]byte, BLOCKLEN)
	for i := 0; i < BLOCKLEN; i++ {
		key[i] = plaintext[i] ^ plaintext[2*BLOCKLEN + i]
	}
	return key, nil
}

func main() {
	textToEncrypt := "supder duper uber secret key that needs encrypting or elsevery bad horrible things will happen"

	testEncrypt := make([]byte, BLOCKLEN)
	for i := 0; i < BLOCKLEN; i++ {
		testEncrypt[i] = BLOCKLEN
	}
	
	CIPHERTEXT = EncryptText(string(textToEncrypt))

	foundKey, err := AttackerFindKey(CIPHERTEXT)
	if err == nil {
		fmt.Println("found key  = ",foundKey)
		fmt.Println("actual key = ", KEY)
	}
}
