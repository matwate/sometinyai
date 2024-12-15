package sometinyai

import (
	"fmt"
	"math"
	"sort"
	"sync"
)

const MUTATION_COUNT = 2

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

func NewPopulation(size int, inputs int, outputs int) Population {
	p := make(Population, size)
	for i := range p {
		p[i].Genome = NewGenome(inputs, outputs)
	}
	return p
}

func NewSimulation(
	size int,
	inputs int,
	outputs int,
	thresh float64,
	brek ThresholdBreak,
) Simulation {
	return Simulation{
		Population:     NewPopulation(size, inputs, outputs),
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
