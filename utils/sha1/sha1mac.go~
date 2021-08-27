package sha1

func MAC(key []byte, text []byte) [Size]byte {
	input := make([]byte, len(key)); copy(input, key)
	input = append(input, text...)
	return Sum(input)
}
