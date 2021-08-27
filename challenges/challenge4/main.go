package main

import(
	"os"
	"bufio"
	"log"
	"fmt"
	"encoding/hex"
	"math"
)



// https://en.wikipedia.org/wiki/Letter_frequency
var englishLetterHist = map[byte]float64 {
	'a': 8.2,
	'b': 1.5,
	'c': 2.8,
	'd': 4.3,
	'e': 13.0,
	'f': 2.2,
	'g': 2.0,
	'h': 6.1,
	'i': 7.0,
	'j': 0.15,
	'k': 0.77,
	'l': 4.0,
	'm': 2.4,
	'n': 6.7,
	'o': 7.5,
	'p': 1.9,
	'q': 0.095,
	'r': 6.0,
	's': 6.3,
	't': 9.1,
	'u': 2.8,
	'v': 0.98,
	'w': 2.4,
	'x': 0.15,
	'y': 2.0,
	'z': 0.074,
	'A': 8.2,
	'B': 1.5,
	'C': 2.8,
	'D': 4.3,
	'E': 13.0,
	'F': 2.2,
	'G': 2.0,
	'H': 6.1,
	'I': 7.0,
	'J': 0.15,
	'K': 0.77,
	'L': 4.0,
	'M': 2.4,
	'N': 6.7,
	'O': 7.5,
	'P': 1.9,
	'Q': 0.095,
	'R': 6.0,
	'S': 6.3,
	'T': 9.1,
	'U': 2.8,
	'V': 0.98,
	'W': 2.4,
	'X': 0.15,
	'Y': 2.0,
	'Z': 0.074,
}


type TextScore struct {
	text string
	score float64
}


// score : Manhattan distance
func score_text(tst_txt string) float64 {
	strLen := len(tst_txt)
	score := 0.0
	freq := make(map[byte]float64)
	for i := 0; i < strLen; i++{
		_, ok := englishLetterHist[tst_txt[i]]
		if ok {
			_, ok = freq[tst_txt[i]]
			if ok {
				freq[tst_txt[i]] += 1.0
			} else{
				freq[tst_txt[i]] = 1.0
			}
		} else if tst_txt[i] != ' ' {
			//return math.MaxFloat64
			//score += 100.0/float64(strLen)
			score += 1000.0/float64(strLen)

		}
	}
	
	for char, fr := range freq {
		h := (fr/float64(strLen)) * 100.0
		score += (h - englishLetterHist[char]) * (h - englishLetterHist[char])
		//score += math.Abs(h - englishLetterHist[char])
	}

	return score
}


func main(){
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	
	smallestScoreOverall := math.MaxFloat64
	var ans string
	var key byte


	
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		/*fmt.Println(scanner.Text())
		fmt.Println()*/

		cipher_text, err := hex.DecodeString(scanner.Text())
		if err != nil {
		
		}


		test_text := make([]byte, len(cipher_text))
		
		smallestScore := math.MaxFloat64
		var ansCurr string
		var keyCurr byte
		for key_byte := byte(0); key_byte <= /*255*/254; key_byte++ {	
			for i, _ := range cipher_text {
				test_text[i] = cipher_text[i] ^ key_byte
			}
			scr := score_text(string(test_text))


			if smallestScore > scr {
				ansCurr = string(test_text)
				smallestScore = scr
				keyCurr = key_byte
			}
			
		}
		if smallestScore < smallestScoreOverall {
			ans = ansCurr
			smallestScoreOverall = smallestScore
			key = keyCurr
		}
		/*fmt.Printf("key: %d \n", key)
		fmt.Printf("ciphertext: %s \n", ans)*/
	}


	fmt.Printf("key: %d \n", key)
	fmt.Printf("ciphertext: %s \n", ans)
	
}
