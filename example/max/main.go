package main

import (
	"log"
	"math"

	"github.com/mohanson/doa"
	"github.com/mohanson/ga"
)

// Find the maximum value of the function:
// F(x) = sin(10 * x) * x + cos(2 * x) * x, x in [0, 5]

func f(x float64) float64 {
	doa.Doa2(x >= 0, x <= 5)
	return math.Sin(10*x)*x + math.Cos(2*x)*x
}

func main() {
	gas := ga.GAs{
		Option: ga.GAsOption{
			GenemoSize: 10,
			PopSize:    80,
			MaxIter:    200,
			PC:         0.5,
			PM:         0.005,
			Fitness: func(g *ga.Genemo) float64 {
				var c uint64 = 0
				for i := 0; i < 10; i++ {
					a := g.Locus[i] & 1
					c |= a << i
				}
				x := float64(ga.GraycodeDecode(c)) / 1023 * 5
				return f(x)
			},
			Trigger: func(g *ga.GAs) {
				if g.Generation%10 == 0 {
					log.Println("Generation", g.Generation)
					i := ga.FindArgMax(g.Fitness)
					log.Println("Individual", g.Option.Fitness(g.Population[i]))
				}
			},
		},
	}
	gas.Run()
}
