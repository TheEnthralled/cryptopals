package main

import (
	"cryptopals/utils/pkcs7"
	//"fmt"
	"log"
)

func main() {
	goodText := "ICE ICE BABY\x04\x04\x04\x04"
	_, err := pkcs7.ValidatePadding([]byte(goodText))
	if err != nil {
		log.Fatal(err)
	}
	badText := "ICE ICE BABY\x05\x05\x05\x05"
	_, err = pkcs7.ValidatePadding([]byte(badText))
	if err != nil {
		log.Print(err)
	}
	badText = "ICE ICE BABY\x01\x02\x03\x04"
	_, err = pkcs7.ValidatePadding([]byte(badText))
	if err != nil {
		log.Print(err)
	}
}
