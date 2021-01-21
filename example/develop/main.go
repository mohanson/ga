package main

import (
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/mohanson/ga"
)

// Genome is all genetic material of an organism.
type Genemo struct {
	Locus []uint64
}

// Copy Genome.
func (g *Genemo) Copy() *Genemo {
	d := make([]uint64, g.Size())
	copy(d, g.Locus)
	return &Genemo{d}
}

// Size returns size of genemo.
func (g *Genemo) Size() int {
	return len(g.Locus)
}

// GAsOption.
type GAsOption struct {
	GenemoSize int                   // Size of genemo
	PopSize    int                   // Population size
	MaxIter    int                   // Number of evolutionary iterations
	PC         float64               // Crossover rate
	PM         float64               // Mutation rate
	Fitness    func(*Genemo) float64 // Fitness function
	Trigger    func(*GAs)            // Called every iteration
}

// GAs.
type GAs struct {
	Option     GAsOption
	Generation int
	Population []*Genemo
}

func (g *GAs) Run() {
	rand.Seed(time.Now().UnixNano())

	// The population size depends on the nature of the problem, but typically contains several hundreds or thousands of
	// possible solutions. Often, the initial population is generated randomly, allowing the entire range of possible
	// solutions (the search space). Occasionally, the solutions may be "seeded" in areas where optimal solutions are
	// likely to be found.
	g.Population = make([]*Genemo, g.Option.PopSize)
	for i := 0; i < g.Option.PopSize; i++ {
		locus := make([]uint64, g.Option.GenemoSize)
		for j := 0; j < g.Option.GenemoSize; j++ {
			locus[j] = rand.Uint64()
		}
		g.Population[i] = &Genemo{Locus: locus}
	}

	allFitness := make([]float64, g.Option.PopSize)
	for ; g.Generation < g.Option.MaxIter; g.Generation++ {
		g.Option.Trigger(g)

		minFitness := math.MaxFloat64
		for i := 0; i < g.Option.PopSize; i++ {
			f := g.Option.Fitness(g.Population[i])
			allFitness[i] = f
			if f < minFitness {
				minFitness = f
			}
		}
		for i := 0; i < g.Option.PopSize; i++ {
			allFitness[i] = allFitness[i] - minFitness
		}
		cntFitness := 0.0
		for i := 0; i < g.Option.PopSize; i++ {
			cntFitness += allFitness[i]
		}
		chdPop := make([]*Genemo, g.Option.PopSize)
		for i := 0; i < g.Option.PopSize; i++ {
			dp := rand.Float64() * cntFitness
			for j := 0; j < g.Option.PopSize; j++ {
				if dp <= allFitness[j] {
					chdPop[i] = g.Population[j].Copy()
					break
				} else {
					dp -= allFitness[j]
				}
			}
		}

		for i := 0; i < g.Option.PopSize/2; i++ {
			if rand.Float64() < g.Option.PC {
				jcd := rand.Int()%(g.Option.GenemoSize) + 1
				for j := 0; j < jcd; j++ {
					a := chdPop[2*i].Locus[j]
					chdPop[2*i].Locus[j] = chdPop[2*i+1].Locus[j]
					chdPop[2*i+1].Locus[j] = a
				}
			}
		}

		for i := 0; i < g.Option.PopSize; i++ {
			for j := 0; j < g.Option.GenemoSize; j++ {
				if rand.Float64() < g.Option.PM {
					chdPop[i].Locus[j] = rand.Uint64()
				}
			}
		}

		copy(g.Population, chdPop)

	}
}

func main() {
	gas := GAs{
		Option: GAsOption{
			GenemoSize: 10,
			PopSize:    80,
			MaxIter:    200,
			PC:         0.5,
			PM:         0.005,
			Fitness: func(g *Genemo) float64 {
				var c uint32 = 0
				for i := 0; i < 10; i++ {
					a := g.Locus[i] & 1
					c |= uint32(a) << i
				}
				f := float64(ga.GraycodeDecode(c)) / 1023 * 5
				return math.Sin(10*f)*f + math.Cos(2*f)*f
			},
			Trigger: func(g *GAs) {
				log.Println("Generation", g.Generation)
			},
		},
	}
	gas.Run()

	max := 0.0
	for i := 0; i < gas.Option.PopSize; i++ {
		n := gas.Option.Fitness(gas.Population[i])
		if n > max {
			max = n
		}
	}
	log.Println(max)

}
