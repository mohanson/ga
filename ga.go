package ga

import (
	"log"
	"math"
	"math/rand"

	"github.com/godump/doa"
)

// Genome is all genetic material of an organism.
type Genemo struct {
	Locus []uint64
}

// Copy Genome.
func (g Genemo) Copy() Genemo {
	d := make([]uint64, g.Size())
	copy(d, g.Locus)
	return Genemo{d}
}

// Size returns size of genemo.
func (g Genemo) Size() int {
	return len(g.Locus)
}

// GAsOption.
type GAsOption struct {
	// Basic operating parameters. Genemo is an array of uint64.
	GenemoSize int
	// Basic operating parameters. Population size, usually from 20 to 100.
	PopSize int
	// Basic operating parameters. Number of evolutionary iterations, usually from 100 to 500.
	MaxIter int
	// Basic operating parameters. Crossover rate, usually from 0.4 to 0.99.
	PC float64
	// Basic operating parameters. Mutation rate, usually from 0.0001 to 0.1.
	PM float64

	// Necessary parameters. Fitness.
	Fitness func(*Genemo) float64

	// Selection operator.
	SelectionOperator int
	// The optimal retention strategy refers to the selection, crossover and mutation of the best individuals in the
	// group directly into the next generation to avoid the loss of outstanding individuals.
	MaintainTheBestIndividual bool

	// A callback function is triggered every time a new generation is generated.
	Trigger func(*GAs)
}

// GAs.
type GAs struct {
	Option                GAsOption
	Generation            int
	Population            []Genemo
	Fitness               []float64
	BestIndividual        Genemo
	BestIndividualFieness float64
}

// The population size depends on the nature of the problem, but typically contains several hundreds or thousands of
// possible solutions. Often, the initial population is generated randomly, allowing the entire range of possible
// solutions (the search space). Occasionally, the solutions may be "seeded" in areas where optimal solutions are
// likely to be found.
func GAsInitialize(g *GAs) {
	g.Population = make([]Genemo, g.Option.PopSize)
	g.Fitness = make([]float64, g.Option.PopSize)
	for i := 0; i < g.Option.PopSize; i++ {
		locus := make([]uint64, g.Option.GenemoSize)
		for j := 0; j < g.Option.GenemoSize; j++ {
			locus[j] = rand.Uint64()
		}
		g.Population[i] = Genemo{Locus: locus}
	}
}

// Measure the fitness of each individual.
func GAsFitnessMessure(g *GAs) {
	for i := 0; i < g.Option.PopSize; i++ {
		f := g.Option.Fitness(&g.Population[i])
		g.Fitness[i] = f
	}
	i := FindArgMax(g.Fitness)
	if g.Generation == 0 || g.Fitness[i] > g.BestIndividualFieness {
		g.BestIndividual = g.Population[i].Copy()
		g.BestIndividualFieness = g.Fitness[i]
	}
	if g.Option.MaintainTheBestIndividual {
		i := FindArgMin(g.Fitness)
		g.Population[i] = g.BestIndividual.Copy()
		g.Fitness[i] = g.BestIndividualFieness
	}
}

// During each successive generation, a portion of the existing population is selected to breed a new generation.
// Individual solutions are selected through a fitness-based process, where fitter solutions (as measured by a fitness
// function) are typically more likely to be selected. Certain selection methods rate the fitness of each solution and
// preferentially select the best solutions. Other methods rate only a random sample of the population, as the former
// process may be very time-consuming.
func GAsSelect(g *GAs) {
	doa.Doa(g.Option.SelectionOperator < 3)
	switch g.Option.SelectionOperator {
	case 0:
		GAsSelectRanked(g)
	case 1:
		GAsSelectProportional(g)
	case 2:
		GAsSelectStochastic(g)
	}
}

// The Roulette algorithm, that requires every fitness must be greater than or equal to 0.
func GAsSelectCommonRoulette(g *GAs, f []float64) {
	cntFitness := 0.0
	for i := 0; i < g.Option.PopSize; i++ {
		cntFitness += f[i]
	}
	generation := make([]Genemo, g.Option.PopSize)
	for i := 0; i < g.Option.PopSize; i++ {
		bullet := rand.Float64() * cntFitness
		for j := 0; j < g.Option.PopSize; j++ {
			if bullet > f[j] {
				bullet -= f[j]
			} else {
				generation[i] = g.Population[j].Copy()
				break
			}
		}
	}
	g.Population = generation
}

// The proportional selection algorithm is a replay random sampling algorithm, and its characteristic is that the
// probability of being selected is proportional to the fitness.
func GAsSelectProportional(g *GAs) {
	// Fitness scaling converts the raw fitness scores that are returned by the fitness function to values in a range
	// that is suitable for the selection function. The selection function uses the scaled fitness values to select the
	// parents of the next generation. The selection function assigns a higher probability of selection to individuals
	// with higher scaled values.
	minFitness := g.Fitness[FindArgMin(g.Fitness)]
	relFitness := make([]float64, g.Option.PopSize)
	for i := 0; i < g.Option.PopSize; i++ {
		relFitness[i] = g.Fitness[i] - minFitness
	}
	GAsSelectCommonRoulette(g, relFitness)
}

// Rank, scales the raw scores based on the rank of each individual instead of its score. The rank of an individual is
// its position in the sorted scores: the rank of the most fit individual is 1, the next most fit is 2, and so on. The
// rank scaling function assigns scaled values so that the scaled value of an individual with rank n is proportional
// to 1/sqrt(n).
//
// https://www.mathworks.com/help/gads/fitness-scaling.html
func GAsSelectRanked(g *GAs) {
	relFitness := make([]float64, g.Option.PopSize)
	s := ArgSort(g.Fitness)
	for i := 0; i < g.Option.PopSize; i++ {
		f := 1.0 / math.Sqrt(float64(g.Option.PopSize-i))
		relFitness[s[i]] = f
	}
	GAsSelectCommonRoulette(g, relFitness)
}

// Stochastic tournament model. It randomly selects 2 individuals at a time, and then the individuals with higher
// fitness are passed on to the next generation.
func GAsSelectStochastic(g *GAs) {
	generation := make([]Genemo, g.Option.PopSize)
	for i := 0; i < g.Option.PopSize; i++ {
		a := rand.Int() % g.Option.PopSize
		b := rand.Int() % g.Option.PopSize
		if g.Fitness[a] > g.Fitness[b] {
			generation[i] = g.Population[a].Copy()
		} else {
			generation[i] = g.Population[b].Copy()
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
	doa.Doa(g.Option.PopSize&1 == 0)
	if g.Option.Trigger == nil {
		g.Option.Trigger = func(g *GAs) {
			log.Printf("G=%d F=%f\n", g.Generation, g.BestIndividualFieness)
		}
	}
	GAsInitialize(g)
	GAsFitnessMessure(g)
	for ; g.Generation < g.Option.MaxIter; g.Generation++ {
		g.Option.Trigger(g)
		GAsSelect(g)
		GAsCrossover(g)
		GAsMutate(g)
		GAsFitnessMessure(g)
	}
}
