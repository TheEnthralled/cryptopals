package aes

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"testing"
)


/*func TestAES_ECB(t *testing.T) {
	key := "YELLOW SUBMARINE"
	fileContents, err := ioutil.ReadFile("7.txt")
	if err != nil {

	}

	cipherText, err := base64.StdEncoding.DecodeString(string(fileContents))
	if err != nil {
		
	}

	fmt.Println(string(cipherText))	
	out , _ := Decrypt_ECB_PKCS7([]byte(key), cipherText)
	fmt.Println(string(out))

	fmt.Println(string(cipherText))
	out, _ = Parallel_Decrypt_ECB_PKCS7([]byte(key), cipherText)
	fmt.Println(string(out))

  }*/

func TestAES_CTR(t *testing.T) {
	key := "YELLOW SUBMARINE"

	plainText := []byte("Spongeboi me bob! ARGH ARGH ARGH ARGH Bikini Bottom yeah yeah yeah")
	for i := 0; i < 10000; i++ {
		plainText = append(plainText, 'A')
	}


	cipherText := Encrypt_CTR(0, 0, []byte(key), plainText)
	
	
	fmt.Println(string(cipherText))	
	out := Decrypt_CTR(0, 0, []byte(key), cipherText)
	fmt.Println(string(out))

}

func BenchmarkAES_ECB(t *testing.B){
	key := "YELLOW SUBMARINE"
	fileContents, err := ioutil.ReadFile("7.txt")
	if err != nil {

	}

	cipherText, err := base64.StdEncoding.DecodeString(string(fileContents))
	if err != nil {
		
	}

	Decrypt_ECB_PKCS7([]byte(key), cipherText)
}


func BenchmarkParallel_AES_ECB(t *testing.B){
	key := "YELLOW SUBMARINE"
	fileContents, err := ioutil.ReadFile("7.txt")
	if err != nil {

	}

	cipherText, err := base64.StdEncoding.DecodeString(string(fileContents))
	if err != nil {
		
	}

	Parallel_Decrypt_ECB_PKCS7([]byte(key), cipherText)
}
