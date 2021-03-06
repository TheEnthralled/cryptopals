package main

import (
	//"crypto/rand"
	//"bytes"
	"bytes"
	"cryptopals/utils/aes"
	"fmt"
	"log"
	"math/rand"
	"regexp"
	"strings"
)

var randKey string
var randIV string

var prefixStr string = "comment1=cooking%20MCs;userdata="
//var suffixStr string =  ";comment2=%20like%20a%20pound%20of%20bacon"
var suffixStr string =  "comment2=%20like%20a%20pound%20of%20bacon"


func Oracle(yourString string) string {
	re := regexp.MustCompile(`;`)
	yourString = re.ReplaceAllLiteralString(yourString, `";"`)
	re = regexp.MustCompile(`=`)
	yourString = re.ReplaceAllLiteralString(yourString, `"="`)
	
	plainText := append([]byte(prefixStr), []byte(yourString)...)
	plainText =  append(plainText, []byte(suffixStr)...)

	if randKey == "" {
		buff := make([]byte, 16)
		_, err := rand.Read(buff)
		if err != nil {
			log.Print(err)
		}
		randKey = string(buff)

		_, err = rand.Read(buff)
		if err != nil {
			log.Print(err)
		}
		randIV = string(buff)
	}
	return string(aes.Encrypt_CBC_PKCS7([]byte(randKey), []byte(randIV), plainText))
}

func DetectAdminStatus(cipherText string) bool {
	plainText, err := aes.Decrypt_CBC_PKCS7([]byte(randKey), []byte(randIV), []byte(cipherText))
	if err != nil {
		log.Print(err)
		return false
	}

	targetString := ";admin=true;"
	if strings.Contains(string(plainText), targetString) {
		fmt.Println(string(plainText))
		return true
	}
	return false
}

// IMPORTANT QUESTIONS: CAN USER FIELD BE JUST GIBBERISH? AND IF SO, HOW CAN WE RECOVER/PREDICT/MAKE THAT INFO BEFORE CALLING THE ORACLES AND SCRAMBLING THE CIPHERTEXTS?
func DoALittleTrolling() {
	targetString := ";admin=true;"
	myString := "EadminEtrueE"
	xorString := make([]byte, len(targetString))
	for i := 0; i < len(targetString); i++ {
		xorString[i] = targetString[i] ^ myString[i]
	}
	yourString := bytes.Repeat([]byte{'A'}, 128)
	yourString = append(yourString, []byte(myString)...)
	
	cipherText := Oracle(string(yourString))
	for i := 0; i + len(xorString) < len(cipherText); i++ {
		
		newCipherText := make([]byte, len(cipherText))
		copy(newCipherText, []byte(cipherText))
		for j := 0; j < len(xorString); j++ {
			newCipherText[i +j] ^= xorString[j]
		}
		if DetectAdminStatus(string(newCipherText)) {
			fmt.Println("success!")
			return
		}
	}
}

func main(){
	DoALittleTrolling()
}
