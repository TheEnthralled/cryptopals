package mt19937_test

import (
	"testing"
	//"fmt"
	"github.com/seehuhn/mt19937"
	. "cryptopals/utils/mt19937"
	"math/rand"
	"time"
)

func TestMT19937(t * testing.T){
	rng1 := rand.New(New())
	rng2 := rand.New(mt19937.New())

	seed := time.Now().UnixNano()
	rng1.Seed(seed)
	rng2.Seed(seed)

	for i := 0; i < 1e9; i++ {
		x, y := rng1.Int63(), rng2.Int63()
		if x > 0 && y > 0 && x != y {
			t.Error("RNG positive outputs differ")
		}
		//fmt.Printf("%d \n%d \n ---------------- \n", rng1.Int63(), rng2.Int63())
	}
}	
