package main

import(
	"os"
	"encoding/base64"
	"encoding/hex"
	"log"
	"fmt"
)

func main(){
     decode, err := hex.DecodeString(os.Args[1])

     if err != nil {
     	log.Fatal(err)
     }

     encoder := base64.NewEncoder(base64.StdEncoding, os.Stdout)
     encoder.Write(decode)
     encoder.Close()
     fmt.Println()
}
