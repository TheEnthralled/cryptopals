package pkcs7

import (
	"errors"
	"log"
)

func Pad(text string, blockSize int) string {
	bytes := []byte(text)
	padLen := blockSize - (len(bytes) % blockSize)

	if padLen > 256 {
		log.Fatal("illegal pad len")
	}

	paddedString := make([]byte, len(bytes) + padLen)
	copy(paddedString[:len(bytes)], bytes)	
	for i := len(bytes); i < len(paddedString); i++ {
		paddedString[i] = byte(padLen)
	}

	return string(paddedString)
}

func ValidatePadding(text []byte) ([]byte, error) {
	textLen := len(text)
	numPadBytes := int(text[textLen-1])
	for i := textLen - 1; i >= textLen - numPadBytes; i-- {
		if text[i] != byte(numPadBytes) {
			return text, errors.New("pkcs: some pad byte does not equal to the total number of padded bytes")
		}
	}
	return text[0:textLen-numPadBytes], nil
}
