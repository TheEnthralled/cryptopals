package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"
)

func main(){
	decode1, err := hex.DecodeString(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	decode2, err := hex.DecodeString(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}

	// Ensure input hex strings are of same length.
	if len(decode1) != len(decode2) {
		log.Fatal(err)
	}

	res := make([]byte, len(decode1))

	for i := 0; i < len(decode1); i++ {
		res[i] = decode1[i] ^ decode2[i]
	}

	encoder := hex.NewEncoder(os.Stdout)
	encoder.Write(res)
	fmt.Println()	
}
