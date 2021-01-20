package main

import (
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/mohanson/ga"
)

type Gene interface {
	Random()
	Copy() Gene
	Not()
}

type GeneBinary struct {
	X uint8
}

func (g *GeneBinary) Random() {
	g.X = uint8(rand.Int() & 1)
}

func (g *GeneBinary) Copy() Gene {
	return &GeneBinary{
		X: g.X,
	}
}

func (g *GeneBinary) Not() {
	g.X ^= 1
}

func (g *GeneBinary) String() string {
	if g.X == 1 {
		return "1"
	} else {
		return "0"
	}
}

func NewGeneBinary() Gene {
	return &GeneBinary{
		X: uint8(rand.Int() & 1),
	}
}

var (
	_ Gene = (*GeneBinary)(nil)
)

type Genemo struct {
	Gene []Gene
}

func (g *Genemo) Copy() *Genemo {
	gene := make([]Gene, len(g.Gene))
	for i := 0; i < len(g.Gene); i++ {
		gene[i] = g.Gene[i].Copy()
	}
	return &Genemo{Gene: gene}
}

type GAsOption struct {
	GenemoSize int
	PopSize    int
	MaxIter    int
	PC         float64
	PM         float64
	NewGene    func() Gene
	Fitness    func(*Genemo) float64
}

type GAs struct {
	Option     GAsOption
	Pop        []*Genemo
	AbsFitness []float64
	RelFitness []float64
}

func (g *GAs) Run() {
	rand.Seed(time.Now().UnixNano())

	g.Pop = make([]*Genemo, g.Option.PopSize)
	g.AbsFitness = make([]float64, g.Option.PopSize)
	g.RelFitness = make([]float64, g.Option.PopSize)
	for i := 0; i < g.Option.PopSize; i++ {
		genemo := make([]Gene, g.Option.GenemoSize)
		for j := 0; j < g.Option.GenemoSize; j++ {
			genemo[j] = g.Option.NewGene()
		}
		g.Pop[i] = &Genemo{Gene: genemo}
	}
	for t := 0; t < g.Option.MaxIter; t++ {
		log.Println("T", t)
		maxFitness := -math.MaxFloat64
		minFitness := math.MaxFloat64
		for i := 0; i < g.Option.PopSize; i++ {
			f := g.Option.Fitness(g.Pop[i])
			g.AbsFitness[i] = f
			if f > maxFitness {
				maxFitness = f
			}
			if f < minFitness {
				minFitness = f
			}
		}
		log.Println("Max", maxFitness)
		for i := 0; i < g.Option.PopSize; i++ {
			g.RelFitness[i] = g.AbsFitness[i] - minFitness
		}

		cntFitness := 0.0
		for i := 0; i < g.Option.PopSize; i++ {
			cntFitness += g.RelFitness[i]
		}

		chdPop := make([]*Genemo, g.Option.PopSize)
		for i := 0; i < g.Option.PopSize; i++ {
			dp := rand.Float64() * cntFitness
			for j := 0; j < g.Option.PopSize; j++ {
				if dp <= g.RelFitness[j] {
					chdPop[i] = g.Pop[j].Copy()
					break
				} else {
					dp -= g.RelFitness[j]
				}
			}
		}

		for i := 0; i < g.Option.PopSize/2; i++ {
			if rand.Float64() < g.Option.PC {
				jcd := rand.Int()%(g.Option.GenemoSize) + 1
				for j := 0; j < jcd; j++ {
					a := chdPop[2*i].Gene[j].Copy()
					chdPop[2*i].Gene[j] = chdPop[2*i+1].Gene[j].Copy()
					chdPop[2*i+1].Gene[j] = a
				}
			}
		}

		for i := 0; i < g.Option.PopSize; i++ {
			for j := 0; j < g.Option.GenemoSize; j++ {
				if rand.Float64() < g.Option.PM {
					chdPop[i].Gene[j].Random()
				}
			}
		}

		copy(g.Pop, chdPop)

	}
}

// f(x)=sin(10x) * x + cos(2x) * x

func main() {
	gas := GAs{
		Option: GAsOption{
			GenemoSize: 10,
			PopSize:    80,
			MaxIter:    200,
			PC:         0.5,
			PM:         0.005,
			NewGene:    NewGeneBinary,
			Fitness: func(g *Genemo) float64 {
				var c uint32 = 0
				for i := 0; i < 10; i++ {
					a := g.Gene[i].(*GeneBinary)
					c |= uint32(a.X) << i
				}
				f := float64(ga.GraycodeDecode(c)) / 1023 * 5
				return math.Sin(10*f)*f + math.Cos(2*f)*f
			},
		},
	}
	gas.Run()
}
