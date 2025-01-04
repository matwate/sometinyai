package sometinyai

import (
	"fmt"
	"math"
	"sort"
	"sync"

	"github.com/matwate/sometinyai/activation"
)

var MUTATION_COUNT = 2

type Agent struct {
	Genome  *Genome
	Fitness float64
}

type Population []Agent

type Simulation struct {
	Population     Population
	threshold      ThresholdBreak
	thresholdValue float64
}

type ThresholdBreak int

const (
	Highest ThresholdBreak = iota
	Lowest
	Closest
)

func (p Population) Len() int {
	return len(p)
}

func (p Population) Less(i, j int) bool {
	return p[i].Fitness < p[j].Fitness
}

func (p Population) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func NewPopulation(size int, inputs int, outputs int, act func(float64) float64) Population {
	if act == nil {
		act = activation.Relu
	}
	p := make(Population, size)
	for i := range p {
		p[i].Genome = NewGenome(inputs, outputs, act)
	}
	return p
}

func SetMutationCount(count int) {
	MUTATION_COUNT = count
}

func NewSimulation(
	size int,
	inputs int,
	outputs int,
	thresh float64,
	brek ThresholdBreak,
	act func(float64) float64,
) Simulation {
	return Simulation{
		Population:     NewPopulation(size, inputs, outputs, act),
		threshold:      brek,
		thresholdValue: thresh,
	}
}

func (s Simulation) Train(
	max_iter int,
	Fitness func(*Genome) float64,
	breaks ...func(float64) bool,
) Agent {
	p := s.Population
Sim:
	for iter := 0; iter < max_iter; iter++ {

		// Evaluate Fitness concurrently
		var wg sync.WaitGroup

		for i := range p {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				p[i].Fitness = Fitness(p[i].Genome)
			}(i)
		}
		wg.Wait()

		// Sort Population based on Fitness
		switch s.threshold {
		case Highest:
			sort.Slice(p, func(i, j int) bool {
				return p[i].Fitness > p[j].Fitness
			})
		case Lowest:
			sort.Slice(p, func(i, j int) bool {
				return p[i].Fitness < p[j].Fitness
			})
		case Closest:
			sort.Slice(p, func(i, j int) bool {
				return math.Abs(p[i].Fitness-s.thresholdValue) < math.Abs(p[j].Fitness-s.thresholdValue)
			})
		}
		// Keep top performers
		thresh := len(p) / 3
		newPop := make(Population, 0, len(p))

		// Append top performers
		newPop = append(newPop, p[:thresh]...)

		// Generate the rest of the Population
		for i := thresh; i < len(p); i++ {
			parent := newPop[i%thresh]

			// Copy and mutate the Genome
			newGenome := parent.Genome.Copy()
			newGenome.Mutate(MUTATION_COUNT)

			// Create a new agent with the mutated Genome
			newAgent := Agent{
				Genome:  newGenome,
				Fitness: 0,
			}

			// Append the new agent
			newPop = append(newPop, newAgent)
		}

		// Update the Population
		s.Population = newPop
		p = s.Population
		fmt.Println("Iteration: ", iter, "Fitness: ", p[0].Fitness)
		_ = fmt.Sprintf("Iteration: %d Fitness: %f", iter, p[0].Fitness)
		// We can specify a break condition, to signify that the training was successful
		best := p[0].Fitness
		if len(breaks) > 0 {
			for _, b := range breaks {
				if b(best) {
					break Sim
				}
			}
		}
	}

	// Return the best agent
	return p[0]
}

func (s *Simulation) TrainWithMutableArguments(
	total_iter int,
	Fitness func(*Genome, ...interface{}) float64,
	onSuccess func(float64, ...interface{}) ([]interface{}, bool),
	successThreshold float64,
	args ...interface{},
) Agent {
	// Difference between Train and TrainWithMutable arguments is that.
	// 1. The latter can take in additional arguments to the Fitness function
	// 2. The latter requires a function when the training threshold is met, allowing it to modify the arguments on a "succesfull" generation also allowing it to break the training loop
	// onSuccess function should return true if the training should stop, it takes the data you passed in as arguments, along with the best fitness value

	p := s.Population
	if len(args) == 0 {
		fmt.Println("No arguments provided, if this is intended, use Train method instead")
		return s.Train(total_iter, func(g *Genome) float64 {
			return Fitness(g)
		})
	}
Sim:
	for iter := 0; iter < total_iter; iter++ {
		// Evaluate Fitness concurrently
		var wg sync.WaitGroup
		for i := range p {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				p[i].Fitness = Fitness(p[i].Genome, args...)
			}(i)
		}
		wg.Wait()

		// Sort Population based on Fitness
		switch s.threshold {
		case Highest:
			sort.Slice(p, func(i, j int) bool {
				return p[i].Fitness > p[j].Fitness
			})
		case Lowest:
			sort.Slice(p, func(i, j int) bool {
				return p[i].Fitness < p[j].Fitness
			})
		case Closest:
			sort.Slice(p, func(i, j int) bool {
				return math.Abs(p[i].Fitness-s.thresholdValue) < math.Abs(p[j].Fitness-s.thresholdValue)
			})
		}
		// Keep top performers
		thresh := len(p) / 3
		newPop := make(Population, 0, len(p))

		// Append top performers
		newPop = append(newPop, p[:thresh]...)

		// Generate the rest of the Population
		for i := thresh; i < len(p); i++ {
			parent := newPop[i%thresh]

			// Copy and mutate the Genome
			newGenome := parent.Genome.Copy()
			newGenome.Mutate(MUTATION_COUNT)

			// Create a new agent with the mutated Genome
			newAgent := Agent{
				Genome:  newGenome,
				Fitness: 0,
			}

			// Append the new agent
			newPop = append(newPop, newAgent)
		}

		// Update the Population
		s.Population = newPop
		p = s.Population
		fmt.Println("Iteration: ", iter, "Fitness: ", p[0].Fitness, "Data", args)
		_ = fmt.Sprintf("Iteration: %d Fitness: %f", iter, p[0].Fitness)
		// We can specify a break condition, to signify that the training was successful
		best := p[0].Fitness
		switch s.threshold {
		case Highest:
			if best >= successThreshold {
				fmt.Println("Success, calling onSuccess")
				newargs, stop := onSuccess(best, args...)
				if stop {
					break Sim
				}
				args = newargs
			}
		case Lowest:
			if best <= successThreshold {
				fmt.Println("Success, calling onSuccess")
				newargs, stop := onSuccess(best, args...)
				if stop {
					break Sim
				}
				args = newargs

			}
		case Closest:
			if math.Abs(best-successThreshold) <= 0.0001 {
				fmt.Println("Success, calling onSuccess")
				newargs, stop := onSuccess(best, args...)
				if stop {
					break Sim
				}
				args = newargs

			}
		}
	}
	return p[0]
}
