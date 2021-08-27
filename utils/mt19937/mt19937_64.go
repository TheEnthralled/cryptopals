package mt19937

// https://en.wikipedia.org/wiki/Mersenne_Twister

import(
	"log"
)

const (
	w64 = 64
	n64 = 312
	m64 = 156
	r64 = 31
	a64 = 0xB5026F5AA96619E9
	u64 = 29
	d64 = 0x5555555555555555
	s64 = 17
	b64 = 0x71D67FFFEDA60000
	t64 = 37
	c64 = 0xFFF7EEE000000000
	l64 = 43
	f64 = 6364136223846793005
	lower_mask_64 = (1 << r64) - 1
	upper_mask_64 = ((1 << w64) - 1) ^ lower_mask_64
	lower_w_64 = (1 << w64) - 1
)

type twister struct {
	mt []uint64
	index int
}

type MT19937_64 struct {
	*twister
}

func New() MT19937_64 {
	return MT19937_64{&twister{make([]uint64, n64), n64+1}}
}

func (tw *twister) Seed(s int64) {
	seed := uint64(s)
	tw.index = n64
	tw.mt[0] = seed
	for i := 1; i < n64; i++ {
		tw.mt[i] = (f64 * (tw.mt[i-1] ^ (tw.mt[i-1] >> (w64-2))) + uint64(i)) & lower_w_64
	}
}


func (tw *twister) Int63() int64 {
	if tw.index >= n64 {
		if tw.index > n64 {
			log.Fatal("mt19937: must seed the RNG first.")
		}
		tw.twist()
	}

	y := tw.mt[tw.index]
	y = y ^ ((y >> u64) & d64)
	y = y ^ ((y << s64) & b64)
	y = y ^ ((y << t64) & c64)
	y = y ^ (y >> l64)

	tw.index++
	return int64(y)
}



func (tw *twister) twist() {
	for i := 0; i < n64; i++ {
		x := (tw.mt[i] & upper_mask_64) | (tw.mt[(i+1) % n64] & lower_mask_64)
		xA := x >> 1
		if x % 2 != 0 {
			xA = xA ^ a64
		}
		tw.mt[i] = tw.mt[(i+m64) % n64] ^ xA
	}
	tw.index = 0
}


