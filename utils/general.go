package utils

import (
	"crypto/rand"
	"math/big"
	"log"
)


// Retunrns random integer in range [start, end).
func RandInt(start, end int) (int, error) {
	diff := int64(end - start)
	
	j, err := rand.Int(rand.Reader, new(big.Int).SetInt64(diff))
	if err != nil {
		log.Fatal(err)
		return 0, err
	}

	return int(j.Uint64()) + start, nil	
}