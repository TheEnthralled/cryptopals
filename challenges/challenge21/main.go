package main

import(
	"cryptopals/utils/mt19937"
	"fmt"
	"time"
)


func main() {
	rng1 := mt19937.NewMT19937_32()
	seed := uint32(time.Now().UnixNano())
	rng1.Seed(seed)

	for i := 0; i < 10; i++ {
		fmt.Println(rng1.Uint32())
	}
}
