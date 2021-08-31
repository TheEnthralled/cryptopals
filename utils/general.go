package utils

import (
	"crypto/rand"
	"io"
	"log"
	"math/big"
	"net"
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

// Read all data from a net.Conn conenction.
func ReadAll(conn net.Conn) ([]byte, error) {	
	var ret []byte
	for {
		temp_buff := make([]byte, 16)
		_, err := conn.Read(temp_buff)
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			break
		}
		ret = append(ret, temp_buff...)
	}
	return ret, nil
}
