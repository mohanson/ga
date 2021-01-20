package ga

func GraycodeEncode(x uint32) uint32 {
	return x ^ (x >> 1)
}

func GraycodeDecode(x uint32) uint32 {
	x ^= x >> 16
	x ^= x >> 8
	x ^= x >> 4
	x ^= x >> 2
	x ^= x >> 1
	return x
}
