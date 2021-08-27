package main

import(
	"os"
	"bufio"
	"log"
	"fmt"
	"encoding/hex"
)

func AES_ECB_Oracle(ciperText []byte, blockSize int) bool {
	blockFreq := make(map[string]int)
	
	for i := 0; i < len(ciperText); i += blockSize {
		blockFreq[string(ciperText[i : i + blockSize])]++
	}
	
	for _, value := range blockFreq {
		if value > 1 {
			return true
		}
	}
		
	return false
}

func main(){
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text, err := hex.DecodeString(scanner.Text())
		if err != nil {
			log.Fatal(err)
		}

		if AES_ECB_Oracle(text, 16) {
			fmt.Printf("%s : cipher text has been encrypted in ECB mode \n", hex.EncodeToString(text))
		}
	}

	
	
}
