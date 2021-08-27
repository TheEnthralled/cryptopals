package main

import (
	"crypto/rand"
	"cryptopals/utils"
	"cryptopals/utils/aes"
	"errors"
	"fmt"
	"log"

)

func AES_Encryption_Oracle(plainText string, blockSize int, appendRange []int) ([]byte, error) {
	fmt.Println("before append: ", []byte(plainText))

	if len(appendRange) != 2 {
		return nil, errors.New("error: appendRange takes only 2 values.")
	}
	
	key := make([]byte, blockSize)
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}


	appendCount1, err := utils.RandInt(appendRange[0], appendRange[1])	
	appendCount2, err := utils.RandInt(appendRange[0], appendRange[1])
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	
	newPlainText := make([]byte, appendCount1 + len(plainText) + appendCount2)
	_, err = rand.Read(newPlainText[:appendCount1])
	if err != nil {
		return nil, err
	}
	copy(newPlainText[appendCount1:], []byte(plainText))
	_, err = rand.Read(newPlainText[appendCount1 + len(plainText):])
	if err != nil {
		return nil, err
	}

	fmt.Println("after appends: ", newPlainText)
	
	
	mode, err := utils.RandInt(0, 2)
	if err != nil {
		return nil, err
	}
	
	// Encrypt in ECB mode if mode == 0. Else, use CBC mode.
	if mode == 0 {
		fmt.Println("Using ECB mode")
		return aes.Encrypt_ECB_PKCS7(key, newPlainText), nil
	}
	fmt.Println("Using CBC mode")
	iv := make([]byte, blockSize)
	_, err = rand.Read(iv)
	if err != nil {
		return nil, err
	}
	return aes.Encrypt_CBC_PKCS7(key, iv, newPlainText), nil	
}

func Detect_AES_Cipher_Mode() string{
	yourInput := make([]byte, 1024)
	blockSize := 16

	cipherText, _ := AES_Encryption_Oracle(string(yourInput), blockSize, []int{5,11})

	// DO I HAVE TO GUESS THE BLOCK SIZE ALSO?
	
	//maxBlockSize := 64
	//var guessedBlockSize int
	var cipherBlockCounts map[string]int
		
	cipherBlockCounts = make(map[string]int)
	
	for i := 0; i < len(cipherText); i += blockSize {
		blck := string(cipherText[i:i+blockSize])
		cipherBlockCounts[blck]++
		if cipherBlockCounts[blck] > 1 {
			return "AES-ECB"
		}	
	}
	return "AES-CBC"
}


func main() {
	fmt.Println(Detect_AES_Cipher_Mode())
}
