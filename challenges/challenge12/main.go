package main

import (
	"crypto/rand"
	. "cryptopals/utils/aes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	//"cryptopals/utils/pkcs7"
)

var randKey string = ""
var unknownString string = ""

func Oracle(yourString string) []byte {
	if unknownString == "" {
		fileContents, err := ioutil.ReadFile(os.Args[1])
		if err != nil {

		}
		unStr, err  := base64.StdEncoding.DecodeString(string(fileContents))
		if err != nil {

		}
		unknownString = string(unStr)
	}
	if randKey == "" {
		key := make([]byte, 16)
		_, err := rand.Read(key)
		if err != nil {

		}
		randKey = string(key)
	}

	plainText := append([]byte(yourString), []byte(unknownString)...)

	return Encrypt_ECB_PKCS7([]byte(randKey), plainText)
}

// TODO: Make more effificent? Only 3 key sizes for AES (128, 192, 256 bits) and 1 block size (128 bits) though Rijndael cipher family more expansive.
func Detect_AES_ECB() bool {
	maxBlockSize := 64
	minBlockSize := 16
	
	yourString := make([]byte, 3*maxBlockSize)
	for i := 0; i < len(yourString); i++ {
		yourString[i] = 'A'
	}
	cipherText := Oracle(string(yourString))
	
	for blockSize := minBlockSize; blockSize <= maxBlockSize; blockSize++ {	
		seenBlocks := make(map[string]bool)
		for i := 0; i < len(cipherText); i += blockSize {
			if seenBlocks[string(cipherText[i : i + blockSize])] {
				return true
			}
			seenBlocks[string(cipherText[i : i + blockSize])] = true
		}
	}
	return false
}

func Find_AES_ECB_BlockSize() int {
	maxBlockSize := 64 //bytes
	
	yourString := make([]byte, 0, maxBlockSize)
	origCipherText := Oracle(string(yourString))
	var blockSize int
	
	
	for i := 1; i <= maxBlockSize; i++ {
		yourString = append(yourString, 'A')
		cipherText := Oracle(string(yourString))
		cipherTextsMatch := true
		for j := 0; j < len(origCipherText); j++ {
			if origCipherText[j] != cipherText[i + j] {
				cipherTextsMatch = false
				break
			}
		}
		if cipherTextsMatch {
			blockSize = i
			break
		}
	}
	return blockSize
}

// WARNING: DOES NOT DO ANY PADDING VALIDATION
func Find_AES_ECB_UnknownStr(blockSize int) string{
	var plainTextFound bool
	var plainText []byte	

	yourBytes := make([]byte, blockSize)
	for i := 0; i < len(yourBytes); i++ {
		yourBytes[i] = 'A'
	}

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
				possibleTexts[string(Oracle(string(possibleBlock))[:blockSize])] = string(possibleBlock)
			}
			newCipherBlock := []byte(Oracle(string(yourBytes[i+1:])))[len(plainText)-blockSize:len(plainText)]
			str, in := possibleTexts[string(newCipherBlock)]
			if  !in {
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
		/*_, err := pkcs7.ValidatePadding(plainText)
		if err == nil {
			break
		}*/
	}

	//ret, _ := pkcs7.ValidatePadding(plainText[blockSize:])
	//return string(ret)
	return string(plainText[blockSize:])
	
}

func Solution(){
	blockSize := Find_AES_ECB_BlockSize()
	fmt.Println(Detect_AES_ECB())
	fmt.Println(Find_AES_ECB_UnknownStr(blockSize))
}

func main(){
	Solution()
}
