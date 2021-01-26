package ga

import (
	"testing"
)

func TestArgSort(t *testing.T) {
	var x []float64
	var r []int

	x = []float64{}
	r = ArgSort(x)
	if len(r) != 0 {
		t.Fail()
	}

	x = []float64{2.0, 1.0, 1.5, 1.25}
	r = ArgSort(x)
	if r[0] != 1 {
		t.Fail()
	}
	if r[1] != 3 {
		t.Fail()
	}
	if r[2] != 2 {
		t.Fail()
	}
	if r[3] != 0 {
		t.Fail()
	}
}
