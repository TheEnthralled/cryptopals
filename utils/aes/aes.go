// TODO: Test & benchark parallel encrypt, decrypt funcs.

package aes

import (
	"crypto/aes"
	"cryptopals/utils/pkcs7"
	"log"
	"encoding/binary"
)

// TODO: throw error is key is not EXACTLY 16 bytes long.

func Encrypt_ECB_PKCS7(key []byte, plainText []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatal(err)
	}

	newPlainText := []byte(pkcs7.Pad(string(plainText), block.BlockSize()))	
	cipherText := make([]byte, len(newPlainText))

	blockSize := block.BlockSize()
	for i := 0; i < len(cipherText); i += blockSize {
		block.Encrypt(cipherText[  i : i + blockSize],
			      newPlainText[i : i + blockSize])
	}

	return cipherText
}

func Decrypt_ECB_PKCS7(key []byte, cipherText []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatal(err)
	}

	plainText := make([]byte, len(cipherText))
	blockSize := block.BlockSize()
	for i := 0; i < len(cipherText); i += blockSize {
		block.Decrypt(plainText[ i : i + blockSize],
			      cipherText[i : i + blockSize])
	}

	return pkcs7.ValidatePadding(plainText)
}

func Parallel_Encrypt_ECB_PKCS7(key []byte, plainText []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatal(err)
	}
	blockSize := block.BlockSize()
	
	newPlainText := []byte(pkcs7.Pad(string(plainText), blockSize))
	cipherText := make([]byte, len(newPlainText))
	numBlocks := len(cipherText)/blockSize

	var cnt chan int	
	for i := 0; i < len(cipherText); i += blockSize {
		go func(j int){
			blck, _ := aes.NewCipher(key)
			blck.Encrypt(cipherText[  j : j + blockSize],
				     newPlainText[j : j + blockSize])
			cnt <- i/blockSize
		}(i)	
	}

	blockCount := 0
	for {
		<-cnt
		blockCount++
		if blockCount == numBlocks {
			break
		}
	}
	
	return cipherText
}


func Parallel_Decrypt_ECB_PKCS7(key []byte, cipherText []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatal(err)
	}
	blockSize := block.BlockSize()
	numBlocks := len(cipherText)/blockSize
	
	plainText := make([]byte, len(cipherText))
	cnt := make(chan int)
	for i := 0; i < len(cipherText); i += blockSize {
		go func(j int){	
			blck, _ := aes.NewCipher(key)
			blck.Decrypt(plainText[ j : j + blockSize],
				     cipherText[j : j + blockSize])
			cnt <- 1
		}(i)
	}
	
	blockCount := 0
	for {
		<-cnt
		blockCount++
		if blockCount == numBlocks {
			break
		}
	}	
	ret, err := pkcs7.ValidatePadding(plainText)
	return ret, err
}


func Encrypt_CBC_PKCS7(key []byte, initVec []byte, plainText []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatal(err)
	}

	blockSize := block.BlockSize()	
	newPlainText := pkcs7.Pad(string(plainText), blockSize)
	tempCipherText := make([]byte, len(newPlainText) + blockSize)

	copy(tempCipherText[:blockSize], initVec)
	
	for i := 0; i < len(newPlainText); i += blockSize {
		tempBlock := make([]byte, blockSize)
		for j := 0; j < blockSize; j++ {
			tempBlock[j] = newPlainText[i + j] ^ tempCipherText[i + j]
		}
		block.Encrypt(tempCipherText[i + blockSize : i + (blockSize << 1)], tempBlock[0 : blockSize])
	}
	
	return tempCipherText[blockSize: ]
}

func Decrypt_CBC_PKCS7(key []byte, initVec []byte, cipherText []byte) ([]byte, error) {
	plainText := make([]byte, len(cipherText))

	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatal(err)
	}
	blockSize := block.BlockSize()

	for i := 0; i < len(plainText); i += blockSize {
		block.Decrypt(plainText[ i : i + blockSize],
			      cipherText[i : i + blockSize])
	}

	for i := 0; i < blockSize; i++ {
		plainText[i] ^= initVec[i]
	}

	for i := blockSize; i < len(plainText); i++ {
		plainText[i] ^= cipherText[i - blockSize]
	}

	ret, err := pkcs7.ValidatePadding(plainText)
	return ret, err
}

func Parallel_Decrypt_CBC_PKCS7(key []byte, initVec []byte, cipherText []byte) ([]byte, error) {
	plainText := make([]byte, len(cipherText))

	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatal(err)
	}
	blockSize := block.BlockSize()
	
	cnt := make(chan int)
	for i := 0; i < len(plainText); i += blockSize {
		go func(j int){	
			blck, _ := aes.NewCipher(key)
			blck.Decrypt(plainText[ j : j + blockSize],
				     cipherText[j : j + blockSize])
			cnt <- 1
		}(i)
	}

	numBlocks := len(cipherText)/blockSize
	blockCount := 0
	for {
		<-cnt
		blockCount++
		if blockCount == numBlocks {
			break
		}
	}


	for i := 0; i < blockSize; i++ {
		plainText[i] ^= initVec[i]
	}

	for i := blockSize; i < len(plainText); i++ {
		plainText[i] ^= cipherText[i - blockSize]
	}


	ret, err := pkcs7.ValidatePadding(plainText)
	return ret, err
}


// IMPORTANT: MAKE SURE NONCE IS ALREADY IN LITTLE ENDIAN!!!
const NONCE = 0

func Encrypt_CTR(nonce uint64, counter uint64, key []byte, plainText []byte) []byte {
	currKeystreamLen := 0
	block, err := aes.NewCipher(key)
	if err != nil {

	}
	blockLen := block.BlockSize()
	keyStreamBlock := make([]byte, blockLen)
	cipherText := make([]byte, len(plainText))
	buff := make([]byte, 16)
	binary.PutUvarint(buff[:8], nonce)

	for {
		binary.PutUvarint(buff[8:], counter)	
		block.Encrypt(keyStreamBlock, buff)

		for i := 0; i < blockLen && i + currKeystreamLen < len(plainText); i++ {
			cipherText[i + currKeystreamLen] = plainText[i + currKeystreamLen] ^ keyStreamBlock[i]
		}
		
		currKeystreamLen += blockLen
		if currKeystreamLen >= len(plainText) {
			break
		}
		counter++
	}

	return cipherText
}

func Decrypt_CTR(nonce uint64, counter uint64, key []byte, cipherText []byte) []byte {
	return Encrypt_CTR(nonce, counter, key, cipherText)
}
