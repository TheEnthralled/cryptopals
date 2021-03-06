package main

import (
	"crypto/rand"
	"cryptopals/utils/sha1"
	"encoding/binary"
	//"encoding/hex"
	"fmt"
	//"io"
	"log"
)

const SHA1_BLOCK_SIZE = 512
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
	pad = append(pad, make([]byte, 8)...)
	binary.BigEndian.PutUint64(pad[len(pad)-8:], l<<3)
	return pad
}


func Server(message string, clientMac string) bool {
	if len(KEY) == 0 {
		KEY = make([]byte, KEY_LEN)
		_, err := rand.Read(KEY)
		if err != nil {
			log.Fatal(err)
		}
	}
	
	serverMac := sha1.MAC(KEY, []byte(message))
	if string(serverMac[:]) == clientMac {
		return true
	} 
	return false
}

// Return new SHA1 hash.
func MaliciousClient(origMAC [sha1.Size]byte, origData string, appendData string) {
	// Guess KEY len
	for k := 1; k <= 256; k++ {
		text := make([]byte, len([]byte(origData))); copy(text, []byte(origData))
		prefix := make([]byte, k)
		padding := MD_Pad(append(prefix, text...))
		text = append(text, padding...)
		text = append(text, []byte(appendData)...)

		d := sha1.MakeFromHash(origMAC, uint64(k + len(origData) + len(padding)))
		d.Write([]byte(appendData))
		newMAC := d.CheckSum()
	
		if Server(string(text), string(newMAC[:])) {
			fmt.Println("ATTACKER GRANTED ACCESS")
			fmt.Printf("message:     %s \n", string(text))
			fmt.Printf("message MAC: % x \n", newMAC)
		}
	}
}
// Knows secret KEY.
// Returns message, message's MAC
func Client() ([]byte, [sha1.Size]byte){
	const client_message = "comment1=cooking%20mcs;userdata=foo;comment2=%20like%20a%20pound%20of%20bacon"

	if len(KEY) == 0 {
		KEY = make([]byte, KEY_LEN)
		_, err := rand.Read(KEY)
		if err != nil {
			log.Fatal(err)
		}
	}
	clientmac := sha1.MAC(KEY, []byte(client_message))
	if Server(client_message, string(clientmac[:])) {
		fmt.Println("client access granted")
	} else {
		fmt.Println("client access denied")
	}
	
	return []byte(client_message), clientmac	
}

func main() {
	/*s, _ := hex.DecodeString("61626364658")
	s2 := MD_Pad(s)
	for i := 0; i < 64; i+=4 {
		if i % 16 == 0 && i != 0{
			fmt.Println()
		}
		fmt.Printf("%s ", hex.EncodeToString(s2[i:i+4]))
	}
	fmt.Println()
	fmt.Println("------------------")

	h := sha1.HashMakeFrom(sha1.Init0, sha1.Init1, sha1.Init2, sha1.Init3, sha1.Init4, 0)
	io.WriteString(h, "His money is twice tainted:")
	io.WriteString(h, " 'taint yours and 'taint mine.")
	fmt.Printf("% x \n", h.Sum(nil))

	d := sha1.MakeFromHashState(sha1.Init0, sha1.Init1, sha1.Init2, sha1.Init3, sha1.Init4, 0)
	d.Write([]byte("His money is twice tainted:"))
	d.Write([]byte(" 'taint yours and 'taint mine."))
	fmt.Printf("% x \n" ,d.CheckSum())*/
	message, MAC := Client()
	fmt.Printf("message:     %s \nmessage MAC: % x \n", message, MAC)
	MaliciousClient(MAC, string(message), ";admin=true")

}
