package ga

import (
	"math"
	"math/rand"
	"sort"
	"time"

	"github.com/mohanson/doa"
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
	GenemoSize     int                   // Size of genemo
	PopSize        int                   // Population size, usually in [20, 100]
	MaxIter        int                   // Number of evolutionary iterations, usually in [100, 500]
	PC             float64               // Crossover rate, usually in [0.4, 0.99]
	PM             float64               // Mutation rate, usually in [0.0001, 0.1]
	Fitness        func(*Genemo) float64 // Fitness function
	FitnessScaling int                   // Fitness scaling function
	Trigger        func(*GAs)            // Called every iteration
}

// GAs.
type GAs struct {
	Option     GAsOption
	Generation int
	Population []*Genemo
	Fitness    []float64
}

// The population size depends on the nature of the problem, but typically contains several hundreds or thousands of
// possible solutions. Often, the initial population is generated randomly, allowing the entire range of possible
// solutions (the search space). Occasionally, the solutions may be "seeded" in areas where optimal solutions are
// likely to be found.
func GAsInitialize(g *GAs) {
	g.Population = make([]*Genemo, g.Option.PopSize)
	g.Fitness = make([]float64, g.Option.PopSize)
	for i := 0; i < g.Option.PopSize; i++ {
		locus := make([]uint64, g.Option.GenemoSize)
		for j := 0; j < g.Option.GenemoSize; j++ {
			locus[j] = rand.Uint64()
		}
		g.Population[i] = &Genemo{Locus: locus}
	}
}

// Measure the fitness of each individual.
func GAsFitnessMessure(g *GAs) {
	for i := 0; i < g.Option.PopSize; i++ {
		f := g.Option.Fitness(g.Population[i])
		g.Fitness[i] = f
	}
}

const (
	FitnessScalingCmin = 0
	FitnessScalingRank = 1
)

// Fitness scaling converts the raw fitness scores that are returned by the fitness function to values in a range that
// is suitable for the selection function. The selection function uses the scaled fitness values to select the parents
// of the next generation. The selection function assigns a higher probability of selection to individuals with higher
// scaled values.
func GAsFitnessScaling(g *GAs) {
	doa.Doa1(g.Option.FitnessScaling < 2)
	switch g.Option.FitnessScaling {
	case FitnessScalingCmin:
		GAsFitnessScalingCmin(g)
	case FitnessScalingRank:
		GAsFitnessScalingRank(g)
	}
}

// The simplest scaling algorithm.
func GAsFitnessScalingCmin(g *GAs) {
	minFitness := g.Fitness[FindArgMin(g.Fitness)]
	for i := 0; i < g.Option.PopSize; i++ {
		g.Fitness[i] -= minFitness
		doa.Doa1(g.Fitness[i] >= 0)
	}
}

// Rank, scales the raw scores based on the rank of each individual instead of its score. The rank of an individual is
// its position in the sorted scores: the rank of the most fit individual is 1, the next most fit is 2, and so on. The
// rank scaling function assigns scaled values so that the scaled value of an individual with rank n is proportional
// to 1/sqrt(n).
//
// https://www.mathworks.com/help/gads/fitness-scaling.html
func GAsFitnessScalingRank(g *GAs) {
	type elem struct {
		i int
		f float64
	}
	l := make([]elem, g.Option.PopSize)
	for i := 0; i < g.Option.PopSize; i++ {
		l[i] = elem{i: i, f: g.Fitness[i]}
	}
	sort.Slice(l, func(i int, j int) bool {
		return l[i].f > l[j].f
	})
	for i := 0; i < g.Option.PopSize; i++ {
		g.Fitness[l[i].i] = 1.0 / math.Sqrt(float64(i+1))
	}
}

// During each successive generation, a portion of the existing population is selected to breed a new generation.
// Individual solutions are selected through a fitness-based process, where fitter solutions (as measured by a fitness
// function) are typically more likely to be selected. Certain selection methods rate the fitness of each solution and
// preferentially select the best solutions. Other methods rate only a random sample of the population, as the former
// process may be very time-consuming.
func GAsSelect(g *GAs) {
	cntFitness := 0.0
	for i := 0; i < g.Option.PopSize; i++ {
		cntFitness += g.Fitness[i]
	}
	generation := make([]*Genemo, g.Option.PopSize)
	for i := 0; i < g.Option.PopSize; i++ {
		bullet := rand.Float64() * cntFitness
		for j := 0; j < g.Option.PopSize; j++ {
			if bullet > g.Fitness[j] {
				bullet -= g.Fitness[j]
			} else {
				generation[i] = g.Population[j].Copy()
				break
			}
		}
	}
	g.Population = generation
}

// The next step is to generate a second generation population of solutions from those selected through a combination
// of genetic operators: crossover (also called recombination), and mutation. For each new solution to be produced, a
// pair of "parent" solutions is selected for breeding from the pool selected previously. By producing a "child"
// solution using the above methods of crossover and mutation, a new solution is created which typically shares many
// of the characteristics of its "parents". New parents are selected for each new child, and the process continues
// until a new population of solutions of appropriate size is generated. Although reproduction methods that are based
// on the use of two parents are more "biology inspired", some research[3][4] suggests that more than two "parents"
// generate higher quality chromosomes.
//
// These processes ultimately result in the next generation population of chromosomes that is different from the
// initial generation. Generally, the average fitness will have increased by this procedure for the population, since
// only the best organisms from the first generation are selected for breeding, along with a small proportion of less
// fit solutions. These less fit solutions ensure genetic diversity within the genetic pool of the parents and
// therefore ensure the genetic diversity of the subsequent generation of children.
func GAsCrossover(g *GAs) {
	for i := 0; i < g.Option.PopSize/2; i++ {
		if rand.Float64() < g.Option.PC {
			a := 2 * i
			b := a + 1
			location := rand.Int()%(g.Option.GenemoSize) + 1
			for j := 0; j < location; j++ {
				g.Population[a].Locus[j] = g.Population[a].Locus[j] ^ g.Population[b].Locus[j]
				g.Population[b].Locus[j] = g.Population[a].Locus[j] ^ g.Population[b].Locus[j]
				g.Population[a].Locus[j] = g.Population[a].Locus[j] ^ g.Population[b].Locus[j]
			}
		}
	}
}

// Randomly mutate genes.
func GAsMutate(g *GAs) {
	for i := 0; i < g.Option.PopSize; i++ {
		for j := 0; j < g.Option.GenemoSize; j++ {
			if rand.Float64() < g.Option.PM {
				g.Population[i].Locus[j] = rand.Uint64()
			}
		}
	}
}

// Start a cruel survival competition.
func (g *GAs) Run() {
	doa.Doa1(g.Option.PopSize&1 == 0)
	rand.Seed(time.Now().UnixNano())
	GAsInitialize(g)
	for ; g.Generation < g.Option.MaxIter; g.Generation++ {
		GAsFitnessMessure(g)
		GAsFitnessScaling(g)
		g.Option.Trigger(g)
		GAsSelect(g)
		GAsCrossover(g)
		GAsMutate(g)
	}
}
