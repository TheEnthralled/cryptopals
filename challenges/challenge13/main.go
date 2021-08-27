package main

import (
	"cryptopals/utils/aes"
	"cryptopals/utils/pkcs7"
	"math/rand"
	"regexp"
	"log"
	"fmt"
	"strings"
	"bytes"
)

var randKey string

func KVParser(text string) string {
	var b strings.Builder
	b.WriteString("{\n")
	
	pairs := strings.Split(text, "&")
	for i := 0; i < len(pairs)-1; i++ {
		pair := strings.Split(pairs[i], "=")
		fmt.Fprintf(&b, "\t%s: '%s',\n", pair[0], pair[1])
	}
	pair := strings.Split(pairs[len(pairs)-1], "=")
	fmt.Fprintf(&b, "\t%s: '%s'\n}", pair[0], pair[1])
	return b.String()
}

func profileFor(email string) string {
	re := regexp.MustCompile(`=`)
	email = re.ReplaceAllLiteralString(email, `"="`)
	re = regexp.MustCompile(`&`)
	email = re.ReplaceAllLiteralString(email, `"&"`)

	return "email=" + email + "&uid=10&role=user"
}

func OracleGive(yourText string) string {
	if randKey == "" {
		buff := make([]byte, 16)
		_, err := rand.Read(buff)
		if err != nil {
			log.Print(err)
		}
		randKey = string(buff)
	}
	return string(aes.Encrypt_ECB_PKCS7([]byte(randKey), []byte(profileFor(yourText))))
}

func OracleParse(yourText string) string {
	plainText, err := aes.Decrypt_ECB_PKCS7([]byte(randKey), []byte(yourText))
	if err != nil {
		log.Print(err)
		return ""
	}
	return string(plainText)
}


func FindCipherBlockSize() int {
	maxBlockSize := 64 //bytes
        minBlockSize := 4 //bytes

	for blckSz := minBlockSize; blckSz <= maxBlockSize; blckSz++ {
		yourString := bytes.Repeat([]byte{'A'}, blckSz << 2)
                cipherText := OracleGive(string(yourString))
		foundRepeatedBlock := false

		for i := 0; i + blckSz < len(cipherText); i++ {
			firstBlock := cipherText[i : i + blckSz]
                        var ri int
                        testi := i + (blckSz << 1)
                        if  testi > len(cipherText) {
                                ri = len(cipherText)
                        } else {
                                ri = testi
                        }
                        secondBlock := cipherText[i + blckSz : ri]

			if bytes.Compare([]byte(firstBlock), []byte(secondBlock)) == 0 {
				foundRepeatedBlock = true
				break
			}
		}

		if foundRepeatedBlock {
			return blckSz
		}	
	}
	return -1
}

// IMPORTANT: ASSUMES UID WILL LAWAYS BE OF SMAE LENGTH
func Forge(blockSize int) {
	prefix := make([]byte, 0, 3*blockSize)
	for true {
		cipherText := OracleGive(string(prefix))
		if len(cipherText) >= 48 && string(cipherText[16:32]) == string(cipherText[32:48]) {
			break
		}
		prefix = append(prefix, 'a')
	}
	

	prefixInp := string(prefix[:len(prefix)-2*blockSize])
	targetBlock := pkcs7.Pad("admin", blockSize)
	cipherText := OracleGive(prefixInp + targetBlock)
	targetBlock = string(cipherText[blockSize:2*blockSize])

	replaceBlock := pkcs7.Pad("user", blockSize)
	cipherText = OracleGive(prefixInp + replaceBlock)
	replaceBlock = string(cipherText[blockSize:2*blockSize])

	prefix = make([]byte, 0, 3*blockSize)
	for true {
		cipherText := OracleGive(string(prefix))
		if len(cipherText) >= blockSize && string(cipherText[len(cipherText)-blockSize:]) == replaceBlock {
			break
		}
		prefix = append(prefix, 'a')
	}
	
	regCipherText := []byte(OracleGive(string(prefix)))
	copy(regCipherText[len(regCipherText)-blockSize:], []byte(targetBlock))
	fmt.Println(KVParser(OracleParse(string(regCipherText))))
}

func main(){
	Forge(FindCipherBlockSize())
}
