package simulation

import (
	"fmt"
	"math"
	"sort"
	"sync"

	"github.com/matwate/sometinyai"
	"github.com/matwate/sometinyai/activation"
)

type (
	Agent struct {
		Genome  *sometinyai.Genome
		Fitness float64
	}
	Population []Agent
	Simulation struct {
		Population Population
		Config     Options
	}
	Options struct {
		population     int
		mutation_count int
		iterations     int
		fitness        func(*sometinyai.Genome, ...interface{}) float64
		breakAbove     bool
		breakBelow     bool
		breakClosest   bool
		breakValue     float64
		closeValue     float64
		useMutableData bool
		mutableData    []interface{}
		dataChange     func(...interface{}) []interface{}
		dataCondition  func(float64, ...interface{}) bool
	}
	Option func(*Options)
)

func PopulationSize(size int) Option {
	return func(o *Options) {
		o.population = size
	}
}

func MutationCount(count int) Option {
	return func(o *Options) {
		o.mutation_count = count
	}
}

func Iterations(iter int) Option {
	return func(o *Options) {
		o.iterations = iter
	}
}

func BreakAbove(value float64) Option {
	return func(o *Options) {
		o.breakAbove = true
		o.breakValue = value
	}
}

func BreakBelow(value float64) Option {
	return func(o *Options) {
		o.breakBelow = true
		o.breakValue = value
	}
}

func BreakClosest(value float64) Option {
	return func(o *Options) {
		o.breakClosest = true
		o.breakValue = value
	}
}

func UseMutableData(data ...interface{}) Option {
	return func(o *Options) {
		o.useMutableData = true
		o.mutableData = data
	}
}

func DataChange(f func(...interface{}) []interface{}) Option {
	return func(o *Options) {
		o.dataChange = f
	}
}

func DataCondition(f func(float64, ...interface{}) bool) Option {
	return func(o *Options) {
		o.dataCondition = f
	}
}

func Fitness(f func(*sometinyai.Genome, ...interface{}) float64) Option {
	return func(o *Options) {
		o.fitness = f
	}
}

func NewSimulation(
	inputs, outputs int,
	activation func(float64) float64,
	opts ...Option,
) Simulation {
	/*k
	  For our default values we will be using the follwing:
	  - Population size of 100
	  - Mutation count of 1
	  - 1000 iterations
	  - No breaking conditions
	  - No mutable data
	*/
	args := Options{
		population:     100,
		mutation_count: 1,
		iterations:     1000,
		fitness:        nil,
		breakAbove:     false,
		breakBelow:     false,
		breakClosest:   false,
		breakValue:     0,
		closeValue:     0,
		useMutableData: false,
		dataChange:     nil,
		dataCondition:  nil,
	}

	for _, opt := range opts {
		opt(&args)
	}

	return Simulation{
		Population: NewPopulation(args.population, inputs, outputs, activation),
		Config:     args,
	}
}

func NewPopulation(size int, inputs int, outputs int, act func(float64) float64) Population {
	if act == nil {
		act = activation.Relu
	}
	p := make(Population, size)
	for i := range p {
		p[i].Genome = sometinyai.NewGenome(inputs, outputs, act)
	}
	return p
}

func (s Simulation) Train() (Agent, []interface{}) {
	p := s.Population
Sim:
	for iter := 0; iter < s.Config.iterations; iter++ {
		var wg sync.WaitGroup
		for i := range p {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				p[i].Fitness = s.Config.fitness(p[i].Genome, s.Config.mutableData...)
			}(i)

		}
		wg.Wait()
		// Sort Population based on Fitness
		if s.Config.breakAbove {
			sort.Slice(p, func(i, j int) bool {
				return p[i].Fitness > p[j].Fitness
			})
		} else if s.Config.breakBelow {
			sort.Slice(p, func(i, j int) bool {
				return p[i].Fitness < p[j].Fitness
			})
		} else if s.Config.breakClosest {
			sort.Slice(p, func(i, j int) bool {
				return math.Abs(p[i].Fitness-s.Config.breakValue) < math.Abs(p[j].Fitness-s.Config.breakValue)
			})
		}

		// Keep top performers
		elite := len(p) / 3
		newPop := make(Population, 0, len(p))

		// Append top performers
		newPop = append(newPop, p[:elite]...)

		for i := elite; i < len(p); i++ {
			parent := newPop[i%elite]
			child := parent.Genome.Copy()
			child.Mutate(s.Config.mutation_count)

			newAgent := Agent{
				Genome:  child,
				Fitness: 0,
			}

			newPop = append(newPop, newAgent)
		}

		s.Population = newPop
		p = s.Population

		best := p[0]

		if s.Config.useMutableData {
			if s.Config.dataCondition(best.Fitness, s.Config.mutableData...) {
				s.Config.mutableData = s.Config.dataChange(s.Config.mutableData...)
			}
		}

		if s.Config.breakAbove {
			if best.Fitness > s.Config.breakValue {
				break Sim
			}
		} else if s.Config.breakBelow {
			if best.Fitness < s.Config.breakValue {
				break Sim
			}
		} else if s.Config.breakClosest {
			if math.Abs(best.Fitness-s.Config.breakValue) < s.Config.closeValue {
				break Sim
			}
		}
		fmt.Println("Iteration:", iter, "Fitness:", best.Fitness)

	}
	return p[0], s.Config.mutableData
}
