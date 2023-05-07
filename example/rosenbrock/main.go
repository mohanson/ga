package main

import (
	"math"

	"github.com/godump/doa"
	"github.com/mohanson/ga"
)

// In mathematical optimization, the Rosenbrock function is a non-convex function, introduced by Howard H. Rosenbrock
// in 1960, which is used as a performance test problem for optimization algorithms.[1] It is also known as
// Rosenbrock's valley or Rosenbrock's banana function.
//
// F(x, y) = 100(x² - y)² + (1 - x)², x, y ∈ [-2.048, 2.048]
// F( 2.048, -2.048) = 3897.7342
// F(-2.048, -2.048) = 3905.9262

func f(x float64, y float64) float64 {
	doa.Doa(math.Abs(x) <= 2.048)
	doa.Doa(math.Abs(y) <= 2.048)
	return 100*math.Pow(math.Pow(x, 2)-y, 2) + math.Pow(1-x, 2)
}

func genemoDecode(g *ga.Genemo) (float64, float64) {
	var x uint64
	for i, e := range g.Locus[00:10] {
		x |= (e & 1) << i
	}
	var y uint64
	for i, e := range g.Locus[10:20] {
		y |= (e & 1) << i
	}
	a := 4.096*float64(ga.GraycodeDecode(x))/1023 - 2.048
	b := 4.096*float64(ga.GraycodeDecode(y))/1023 - 2.048
	return a, b
}

func main() {
	gas := ga.GAs{
		Option: ga.GAsOption{
			GenemoSize: 20,
			PopSize:    80,
			MaxIter:    200,
			PC:         0.6,
			PM:         0.001,
			Fitness: func(g *ga.Genemo) float64 {
				return f(genemoDecode(g))
			},
		},
	}
	gas.Run()
}
