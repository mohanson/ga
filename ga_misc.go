package ga

func GraycodeEncode(x uint64) uint64 {
	return x ^ (x >> 1)
}

func GraycodeDecode(x uint64) uint64 {
	x ^= x >> 32
	x ^= x >> 16
	x ^= x >> 8
	x ^= x >> 4
	x ^= x >> 2
	x ^= x >> 1
	return x
}

func FindArgMax(x []float64) int {
	sx := x[0]
	si := 0
	for i := 0; i < len(x); i++ {
		if x[i] > sx {
			sx = x[i]
			si = i
		}
	}
	return si
}

func FindArgMin(x []float64) int {
	sx := x[0]
	si := 0
	for i := 0; i < len(x); i++ {
		if x[i] < sx {
			sx = x[i]
			si = i
		}
	}
	return si
}

// ArgSort returns the indices that would sort an array.
// https://numpy.org/doc/stable/reference/generated/numpy.argsort.html
func ArgSort(x []float64) []int {
	n := len(x)
	a := make([]float64, n)
	copy(a, x)
	r := make([]int, len(x))
	for i := 0; i < n; i++ {
		r[i] = i
	}
	if n < 2 {
		return r
	}
	for i := 1; i < n; i++ {
		v := a[i]
		b := r[i]
		k := i - 1
		for k >= 0 && v < a[k] {
			a[k+1] = a[k]
			r[k+1] = r[k]
			k--
		}
		a[k+1] = v
		r[k+1] = b
	}
	return r
}
