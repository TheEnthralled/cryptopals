package main

import(
	"fmt"
	"cryptopals/utils/pkcs7"
)


func main(){
	plainText := "YELLOW SUBMARINE"
	fmt.Printf("%q", pkcs7.Pad(plainText, 20))
}
