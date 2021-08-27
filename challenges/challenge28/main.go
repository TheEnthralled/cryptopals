package main

import(
	"cryptopals/utils/sha1"
	"fmt"
)

func main() {
	data := []byte("This page intentionally left blank.")
	fmt.Printf("% x \n", sha1.Sum(data))

	key := "YELLOW SUBMARINE"
	fmt.Printf("% x \n", sha1.MAC([]byte(key), data))

	test := []byte("YELLOW SUBMARINEThis page intentionally left blank.")
	fmt.Printf("% x \n", sha1.Sum(test))

}
