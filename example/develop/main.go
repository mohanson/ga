package main

import (
	"log"
	"math"

	"github.com/mohanson/ga"
)

func main() {
	gas := ga.GAs{
		Option: ga.GAsOption{
			GenemoSize: 10,
			PopSize:    80,
			MaxIter:    200,
			PC:         0.5,
			PM:         0.005,
			Fitness: func(g *ga.Genemo) float64 {
				var c uint32 = 0
				for i := 0; i < 10; i++ {
					a := g.Locus[i] & 1
					c |= uint32(a) << i
				}
				f := float64(ga.GraycodeDecode(c)) / 1023 * 5
				return math.Sin(10*f)*f + math.Cos(2*f)*f
			},
			Trigger: func(g *ga.GAs) {
				log.Println("Generation", g.Generation)
				maxFitness := 0.0
				maxIndividualIndex := 0
				for i := 0; i < g.Option.PopSize; i++ {
					f := g.Fitness[i]
					if f > maxFitness {
						maxFitness = f
						maxIndividualIndex = i
					}
				}
				log.Println(g.Option.Fitness(g.Population[maxIndividualIndex]))
			},
		},
	}
	gas.Run()
}
