package main

import (
	//"cryptopals/utils"
	"cryptopals/utils/mt19937"
	"fmt"
	"time"
	//"cryptopals/utils"
)

func DiscoverSeed(r uint32) uint32 {
	var seed uint32
	for i := uint32(0); i < 0xFFFFFFFF; i++ {
		rng := mt19937.NewMT19937_32()
		rng.Seed(i)
		if rng.Uint32() == r {
			seed = i
			break
		}	
	}
	return seed
}

func main() {
	rng := mt19937.NewMT19937_32()
	
	t, _ := utils.RandInt(40, 1001)
	time.Sleep(time.Duration(t) * time.Second)
	seed := uint32(time.Now().UnixNano())
	rng.Seed(uint32(seed))
	t, _ = utils.RandInt(40, 1001)
	time.Sleep(time.Duration(t) * time.Second)

	out := rng.Uint32()
	fmt.Println(out)
	fmt.Println(seed)
	fmt.Println(DiscoverSeed(out))
}
