package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/bits"
	"os"
	"strconv"

	//"os"
	"encoding/base64"
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
	/*' ': 0.0,
	'\n': 0.0,*/
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
	'\t' : 0.0,
	'\n' : 0.0,
	' ' : 0.0,
	/*'1' : 0.0,
	'2' : 0.0,
	'3' : 0.0,
	'4' : 0.0,
	'5' : 0.0,
	'6' : 0.0,
	'7' : 0.0,
	'8' : 0.0,
	'9' : 0.0,
	'0' : 0.0,*/
}

// score : Manhattan distance
func ScoreText(tst_txt string) float64 {
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
		}else{ 
			//return math.MaxFloat64
			//score += 100.0/float64(strLen)
			score += 100000000000.0/float64(strLen)
		}
	}
	
	for char, fr := range freq {
		h := (fr/float64(strLen)) * 100.0
		//score += (h - englishLetterHist[char]) * (h - englishLetterHist[char])
		score += math.Abs(h - englishLetterHist[char])
	}

	return score
}

// Works (most likely).
func HammingDist(str1, str2 string) int {
	bytes1, bytes2 := []byte(str1), []byte(str2)
	len1, len2 := len(bytes1), len(bytes2)
	var lenShorterStr int
	if len1 < len2 {
		lenShorterStr = len1
	} else {
		lenShorterStr = len2
	}
	
	dist := 0
	for i := 0; i < lenShorterStr; i++ {
		dist += bits.OnesCount(uint(bytes1[i] ^ bytes2[i]))
	}
	
	// Add to dist number of 1 bits remaining in the tail of the longer string.
	if len1 != len2 {
		var longerStr []byte
		if len1 > len2 {
			longerStr = bytes1
		} else {
			longerStr = bytes2
		}
		for i := 0; i < len(longerStr) - lenShorterStr; i++ {
			dist += bits.OnesCount(uint(longerStr[i + lenShorterStr]))
		}
	}

	return dist
}


func BreakRepeatingKeyXOR(ciperText string){
	// Find keyLen
	var keyLen int
	keyMinHammDist := math.MaxFloat64
	cipherBytes := []byte(ciperText)
	
	for guessedKeyLen := 3; guessedKeyLen <= 40; guessedKeyLen++ {
		firstBytes  := make([]byte, guessedKeyLen)
		secondBytes := make([]byte, guessedKeyLen)
		
		for i := 0; i < guessedKeyLen; i++ {
			firstBytes[i] = cipherBytes[i + guessedKeyLen]
			secondBytes[i] = cipherBytes[i + 2*guessedKeyLen]

		}

		normDist := float64(HammingDist(string(firstBytes), string(secondBytes)))/float64(guessedKeyLen)
		
		fmt.Printf("guessedKeyLen = %d, normDist = %f \n", guessedKeyLen,normDist)
		
		if normDist < keyMinHammDist {
			keyMinHammDist = normDist
			keyLen = guessedKeyLen
		}
	}
	fmt.Printf("keyLen = %d \n", keyLen)

	
	transposedBlocks := make([][]byte, keyLen)

	for i := 0; i < len(cipherBytes); i++ {
		transposedBlocks[i % keyLen] = append(transposedBlocks[i % keyLen], cipherBytes[i])
	}

	for i := 0; i < keyLen; i++ {
		fmt.Printf("transposedBlocks[%d]: \n", i)
		fmt.Println(transposedBlocks[i])
		fmt.Println()
	}

	
	key := make([]byte, keyLen)
	for i := 0; i < keyLen; i++ {
		testText := make([]byte, len(transposedBlocks[i]))
		smallestScore := math.MaxFloat64
		var kb byte
		for keyByte := 0; keyByte < 256; keyByte++ {
			for j, _ := range transposedBlocks[i] {
				testText[j] = transposedBlocks[i][j] ^ byte(keyByte)
			}
			scr := ScoreText(string(testText))
			if scr < smallestScore {
				smallestScore = scr
				kb = byte(keyByte)
			}
		}
		key[i] = kb
	}


	decryptedText := make([]byte, len(cipherBytes))
	for i := 0; i < len(cipherBytes); i++ {
		decryptedText[i] = cipherBytes[i] ^ key[i % keyLen]
	}
	fmt.Println(decryptedText)
	fmt.Println(string(decryptedText))

	fmt.Println()
	fmt.Println("key: ")
	fmt.Println(string(key))
	fmt.Println()


}



func main(){
	fmt.Println(HammingDist("this is a test","wokka wokka!!!"))
	fileContents, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(fileContents))
	fmt.Println(strconv.Quote(string(fileContents)))
	fmt.Println("-------------------------")

	cipherStr, err := base64.StdEncoding.DecodeString(string(fileContents))
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Println(cipherStr)
	BreakRepeatingKeyXOR(string(cipherStr))
}
