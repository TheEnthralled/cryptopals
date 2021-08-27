package main

import (
	"crypto/rand"
	"cryptopals/utils/md4"
	//"encoding/binary"
	"fmt"
	"log"
)

var KEY []byte
var KEY_LEN int = 16

func MD_Pad(text []byte) []byte {
	l := uint64(len(text))
	var tmp [64]byte
	tmp[0] = 0x80

	var pad []byte
	if l%64 < 56 {
		pad = tmp[:56-l%64]
	} else {
		pad = tmp[:64+56-l%64]
	}
	final := make([]byte, 8)
	l <<= 3
	for i := uint(0); i < 8; i++ {
		final[i] = byte(l >> (8*i))
	}
	return append(pad, final...)
}

func Server(message string, clientMAC string) bool {
	if len(KEY) == 0 {
		KEY = make([]byte, KEY_LEN)
		_, err := rand.Read(KEY)
		if err != nil {
			log.Fatal(err)
		}
	}

	serverMAC := md4.MAC(KEY, []byte(message))
	if string(serverMAC) == clientMAC {
		return true
	}
	return false
}

func Client() ([]byte, []byte) {
	const client_message = "comment1=cooking%20mcs;userdata=foo;comment2=%20like%20a%20pound%20of%20bacon"

	if len(KEY) == 0 {
		KEY = make([]byte, KEY_LEN)
		_, err := rand.Read(KEY)
		if err != nil {
			log.Fatal(err)
		}
	}
	clientmac := md4.MAC(KEY, []byte(client_message))
	if Server(client_message, string(clientmac[:])) {
		fmt.Println("client access granted")
	} else {
		fmt.Println("client access denied")
	}
	
	return []byte(client_message), clientmac	
}


func MaliciousClient(origMAC []byte, origData string, appendData string) {
	// Guess KEY len
	for k := 1; k <= 256; k++ {
		text := make([]byte, len([]byte(origData))); copy(text, []byte(origData))
		prefix := make([]byte, k)
		padding := MD_Pad(append(prefix, text...))
		text = append(text, padding...)
		text = append(text, []byte(appendData)...)

		d := md4.StartFromHash(origMAC, uint64(k + len(origData) + len(padding)))
		d.Write([]byte(appendData))
		newMAC := d.Sum(nil)
	
		if Server(string(text), string(newMAC)) {
			fmt.Println("ATTACKER GRANTED ACCESS")
			fmt.Printf("message:     %s \n", string(text))
			fmt.Printf("message MAC: % x \n", newMAC)
		}
	}
}


func main() {
	message, MAC := Client()
	fmt.Printf("message:     %s \nmessage MAC: % x \n", message, MAC)
	MaliciousClient(MAC, string(message), ";admin=true")


}
