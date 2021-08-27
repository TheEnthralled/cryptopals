package md4

func MAC(key []byte, text []byte) []byte {
	d := new(Digest)
	input := make([]byte, len(key)); copy(input, key)
	input = append(input, text...)
	d.Write(input)
	return d.Sum(nil)
}
