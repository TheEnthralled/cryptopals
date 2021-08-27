package main

import (
	"crypto/rand"
	"cryptopals/utils"
	"cryptopals/utils/aes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"

	//"math"
	"bytes"
)


var randKey string = ""
var randPrefix string = ""
var unknownString string = ""

// TODO: Test with blocksizes of length other than 16 bytes.
func Oracle(yourString string) string {
	if unknownString == "" {
                fileContents, err := ioutil.ReadFile(os.Args[1])
                if err != nil {

                }
                unStr, err  := base64.StdEncoding.DecodeString(string(fileContents))
                if err != nil {

                }
                unknownString = string(unStr)

		// Generate random key.
		key := make([]byte, 16)
                _, err = rand.Read(key)
                if err != nil {

                }
                randKey = string(key)


		// Generate random prefix.
		prefixLen, err:= utils.RandInt(0, 100)
		if err != nil {

		}
		prefix := make([]byte, prefixLen)
		_, err = rand.Read(prefix)
		if err != nil {

		}
		randPrefix = string(prefix)	
        }
        
	plainText := append([]byte(randPrefix), []byte(yourString)...)
	plainText = append(plainText, []byte(unknownString)...)

	return string(aes.Encrypt_ECB_PKCS7([]byte(randKey), plainText))
}

func Find_BlockCipher_Info() (int, int, []byte) {
	maxBlockSize := 64 //bytes
	minBlockSize := 4 //bytes
	
	for blckSz := minBlockSize; blckSz <= maxBlockSize; blckSz++ {
		yourString := bytes.Repeat([]byte{'A'}, blckSz << 2)
		cipherText := Oracle(string(yourString))
		foundRepeatedBlock := false
		var prefixLenUpperBound int
		var encryptedAllABlock []byte
		for i := 0; i + blckSz < len(cipherText); i+= blckSz {
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
				prefixLenUpperBound = i
				encryptedAllABlock = []byte(firstBlock)
				break
			}
		}
		if foundRepeatedBlock {
			return blckSz, prefixLenUpperBound, encryptedAllABlock
		}
	}
	return -1, -1, nil
}

func Find_PrefixLen(blockSize int, prefixLenUpperBound int, repeatedBlock string) int {
	yourString := make([]byte, 0, 3*blockSize)
	for i := 0; i < cap(yourString); i++ {
		block := Oracle(string(yourString))[prefixLenUpperBound : prefixLenUpperBound + blockSize]
		if block == repeatedBlock {
			return prefixLenUpperBound - (i-blockSize)
		}	
		yourString = append(yourString, 'A')
	}
	return -1
}

func Decrypt_UnknownText(blockSize int, prefixLenUpBound int, prefixLen int) string {
	var plainTextFound bool
	var plainText []byte

	offsetBlockLen := prefixLenUpBound - prefixLen
	offsetBlock := bytes.Repeat([]byte{'A'}, offsetBlockLen)	

	yourBytes := bytes.Repeat([]byte{'A'}, blockSize)

	plainText = append(plainText, yourBytes...)

	for true {
		newPlainTextBlock := make([]byte, 0, blockSize)
		prevPlainTextBlock := plainText[len(plainText)-blockSize:]

		for i := 0; i < blockSize; i++ {
			possibleBlock := make([]byte, 0, blockSize)
                        possibleBlock = append(possibleBlock, prevPlainTextBlock[i+1:]...)
                        possibleBlock = append(possibleBlock, newPlainTextBlock...)
                        possibleBlock = append(possibleBlock, '0')
	
			possibleTexts := make(map[string]string)
			for b := 0; b < (1 << 8); b++ {
				possibleBlock[blockSize-1] = byte(b)
				inputBytes := append(offsetBlock, possibleBlock...)
				cipherText := Oracle(string(inputBytes))
				possibleTexts[cipherText[prefixLenUpBound : prefixLenUpBound + blockSize]] = string(possibleBlock)
			}
			cipherText := Oracle(string(append(offsetBlock, yourBytes[i+1:]...)))
			cipherBlock := cipherText[prefixLenUpBound + len(plainText) - blockSize : prefixLenUpBound + len(plainText)]
			str, in := possibleTexts[cipherBlock]
			if !in {
				plainTextFound = true
				break
			}
			newByte := []byte(str)[blockSize-1]
			newPlainTextBlock = append(newPlainTextBlock, newByte)
		}
		plainText = append(plainText, newPlainTextBlock...)
		if plainTextFound {
			break
		}
	}
	return string(plainText[blockSize:])
}

func main(){
	//fmt.Println(Find_Cipher_Block_Size())
	blockSize, prefixLenUpperBound, repeatedBlock := Find_BlockCipher_Info()
	fmt.Println("prefixLenUpperBound= ", prefixLenUpperBound)
	fmt.Println("actual prefixLen= ", len(randPrefix))
	prefixLen := Find_PrefixLen(blockSize, prefixLenUpperBound, string(repeatedBlock))
	fmt.Println("prefixLen= ", prefixLen)
	fmt.Println(Decrypt_UnknownText(blockSize, prefixLenUpperBound, prefixLen))
}
