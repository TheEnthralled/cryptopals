package main

import(
	"cryptopals/utils/mt19937"
	"fmt"
	"time"
)


const BITLEN = 32

func Submask(a uint32, i int, j int) uint32 {
	return (a >> i) & ((1 << (j-i)) - 1)
}

func InvRightOp(z uint32, offset int) uint32 {
	y := uint32(0)
	y ^= Submask(z, BITLEN-offset, BITLEN) << (BITLEN-offset)
	i := BITLEN-offset; for ; i-offset >= 0; i -= offset {
		y ^= (Submask(z, i-offset, i) ^ Submask(y, i, i+offset)) << (i-offset)
	}
	y ^= Submask(z, 0, i) ^ Submask(y, offset, offset + i)
	return y
}

func InvLeftOp(z uint32, offset int, c uint32) uint32 {
	y := uint32(0)
	y ^= Submask(z, 0, offset)
	i := offset; for ; i+offset <= BITLEN; i += offset {
		y ^= (Submask(z, i, i+offset) ^ (Submask(y, i-offset, i) & Submask(c, i, i+offset))) << i
	}
	y ^= (Submask(z, i, BITLEN) ^ (Submask(y, i-offset, BITLEN-offset) & Submask(c, i, BITLEN))) << i
	return y
}

func Untemper_MT19937_32(z uint32) uint32 {
	y := InvRightOp(z, 18)
	y  = InvLeftOp(y, 15, 0xEFC60000)
	y  = InvLeftOp(y, 7,  0x9D2C5680)
	y  = InvRightOp(y, 11)
	return y
}

func main() {
	rand1Output := make([]uint32, 624)
	initState := make([]uint32, 624)
	seed := uint32(time.Now().UnixNano())

	rng1 := mt19937.NewMT19937_32()
	rng1.Seed(seed)

	for i := 0; i < 624; i++ {
		rand1Output[i] = rng1.Uint32()
		initState[i] = Untemper_MT19937_32(rand1Output[i])
	}

	rng2 := mt19937.MT19937_32{&mt19937.Twister32{initState, 0}}

	numMismatch := 0
	for i := 0; i < 624; i++ {
		out := rng2.Uint32()
		fmt.Printf("%d %d \n", out, rand1Output[i])
		if out != rand1Output[i] {
			numMismatch++
		}
	}
	fmt.Println("num mismatch = ", numMismatch)	
}
