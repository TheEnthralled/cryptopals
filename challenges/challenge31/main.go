package main

import (
	"fmt"
	"io"
	"time"
	"encoding/hex"
	"log"
	"net/http"
)


func TimingAttack() {
	
	respBase := "http://localhost:9000/test?file=foo&signature="

	var signature [32]byte
	
	numTrials := 5
	for i := 0; i < len(signature); i++ {

		slowestResponseTime := int64(0)
		nextByte := 0
		
		for b := 0; b < 256; b++ {
			signature[i] = byte(b)

			resp, err := http.Get(respBase + hex.EncodeToString(signature[:]))
			if err != nil {
				log.Fatal(err)
			}
			if resp.StatusCode == 200 {
				goto end
			}	
			resp.Body.Close()

			
			totalDur := int64(0)
			for j := 0; j < numTrials; j++ {
				start := time.Now()
				
				resp, err = http.Get(respBase + hex.EncodeToString(signature[:]))
				if err != nil {
					log.Fatal(err)
				}
				
				dur := time.Since(start).Milliseconds()
				totalDur += dur
				resp.Body.Close()
			}
			
			

			if totalDur > slowestResponseTime {
				slowestResponseTime = totalDur
				nextByte = b
			}
			
			
		}	
		signature[i] = byte(nextByte)	
		fmt.Println("currSig = ", hex.EncodeToString(signature[:]))
	}

end:;


	
	fmt.Println("signature = ", hex.EncodeToString(signature[:]))
	resp, err := http.Get(respBase + hex.EncodeToString(signature[:]))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("status code = ", resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(body))
}

func main() {
	go Server()
	TimingAttack()
}
