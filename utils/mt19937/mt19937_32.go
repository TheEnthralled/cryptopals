package mt19937

// https://en.wikipedia.org/wiki/Mersenne_Twister

import(
	"log"

)

const (
	w32 = 32
	n32 = 624
	m32 = 397
	r32 = 31
	a32 = 0x9908B0DF
	u32 = 11
	d32 = 0xFFFFFFFF
	s32 = 7
	b32 = 0x9D2C5680
	t32 = 15
	c32 = 0xEFC60000
	l32 = 18
	f32 = 1812433253 
	lower_mask_32 = (1 << r32) - 1
	upper_mask_32 = ((1 << w32) - 1) ^ lower_mask_32
	lower_w_32 = (1 << w32) - 1
)

type Twister32 struct {
	State []uint32
	Index int
}

type MT19937_32 struct {
	*Twister32
}

func NewMT19937_32() MT19937_32 {
	return MT19937_32{&Twister32{make([]uint32, n32), n32+1}}
}

func (tw *Twister32) Seed(seed uint32) {
	tw.Index = n32
	tw.State[0] = seed
	for i := 1; i < n32; i++ {
		tw.State[i] = (f32 * (tw.State[i-1] ^ (tw.State[i-1] >> (w32-2))) + uint32(i)) & lower_w_32
	}
}


func (tw *Twister32) Uint32() uint32 {
	if tw.Index >= n32 {
		if tw.Index > n32 {
			log.Fatal("mt19937: must seed the RNG first.")
		}
		tw.twist()
	}

	y := tw.State[tw.Index]
	y = y ^ ((y >> u32) & d32)
	y = y ^ ((y << s32) & b32)
	y = y ^ ((y << t32) & c32)
	y = y ^ (y >> l32)

	tw.Index++
	return y
}



func (tw *Twister32) twist() {
	for i := 0; i < n32; i++ {
		x := (tw.State[i] & upper_mask_32) | (tw.State[(i+1) % n32] & lower_mask_32)
		xA := x >> 1
		if x % 2 != 0 {
			xA = xA ^ a32
		}
		tw.State[i] = tw.State[(i+m32) % n32] ^ xA
	}
	tw.Index = 0
}



