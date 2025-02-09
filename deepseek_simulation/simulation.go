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
		Config     *Options
	}
	Options struct {
		PopulationSize  int
		MutationCount   int
		Iterations      int
		Fitness         func(*sometinyai.Genome, interface{}) float64 // Now takes mutable data
		Threshold       sometinyai.ThresholdBreak
		ThresholdValue  float64
		MutableData     interface{}
		SuccessCallback func(float64, interface{}) (interface{}, bool)
	}
	Option func(*Options)
)

func PopulationSize(size int) Option {
	return func(o *Options) { o.PopulationSize = size }
}

func MutationCount(count int) Option {
	return func(o *Options) { o.MutationCount = count }
}

func Iterations(iter int) Option {
	return func(o *Options) { o.Iterations = iter }
}

func Threshold(threshold sometinyai.ThresholdBreak, value float64) Option {
	return func(o *Options) {
		o.Threshold = threshold
		o.ThresholdValue = value
	}
}

func UseMutableData(
	initialData interface{},
	successCB func(float64, interface{}) (interface{}, bool),
) Option {
	return func(o *Options) {
		o.MutableData = initialData
		o.SuccessCallback = successCB
	}
}

func Fitness(f func(*sometinyai.Genome, interface{}) float64) Option {
	return func(o *Options) { o.Fitness = f }
}

func NewSimulation(inputs, outputs int, act func(float64) float64, opts ...Option) Simulation {
	options := &Options{
		PopulationSize: 100,
		MutationCount:  2,
		Iterations:     1000,
		Threshold:      sometinyai.Highest,
	}

	for _, opt := range opts {
		opt(options)
	}

	return Simulation{
		Population: newPopulation(options.PopulationSize, inputs, outputs, act),
		Config:     options,
	}
}

func newPopulation(size, inputs, outputs int, act func(float64) float64) Population {
	if act == nil {
		act = activation.Relu
	}
	p := make(Population, size)
	for i := range p {
		p[i].Genome = sometinyai.NewGenome(inputs, outputs, act)
	}
	return p
}

func (s Simulation) Train() (Agent, interface{}) {
	for iter := 0; iter < s.Config.Iterations; iter++ {
		var wg sync.WaitGroup

		// Evaluate fitness with current mutable data
		for i := range s.Population {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				s.Population[i].Fitness = s.Config.Fitness(
					s.Population[i].Genome,
					s.Config.MutableData,
				)
			}(i)
		}
		wg.Wait()

		// Sort population based on threshold
		switch s.Config.Threshold {
		case sometinyai.Highest:
			sort.Slice(s.Population, func(i, j int) bool {
				return s.Population[i].Fitness > s.Population[j].Fitness
			})
		case sometinyai.Lowest:
			sort.Slice(s.Population, func(i, j int) bool {
				return s.Population[i].Fitness < s.Population[j].Fitness
			})
		case sometinyai.Closest:
			sort.Slice(s.Population, func(i, j int) bool {
				return math.Abs(s.Population[i].Fitness-s.Config.ThresholdValue) <
					math.Abs(s.Population[j].Fitness-s.Config.ThresholdValue)
			})
		}

		// Breed new generation
		elite := len(s.Population) / 3
		newPop := append(Population{}, s.Population[:elite]...)
		for i := elite; i < len(s.Population); i++ {
			child := s.Population[i%elite].Genome.Copy()
			child.Mutate(s.Config.MutationCount)
			newPop = append(newPop, Agent{Genome: child})
		}
		s.Population = newPop

		bestFitness := s.Population[0].Fitness

		// Check success condition and update mutable data
		if s.Config.SuccessCallback != nil {
			switch s.Config.Threshold {
			case sometinyai.Highest:
				if bestFitness >= s.Config.ThresholdValue {
					if newData, stop := s.Config.SuccessCallback(bestFitness, s.Config.MutableData); stop {
						return s.Population[0], newData
					} else {
						s.Config.MutableData = newData
					}
				}
			case sometinyai.Lowest:
				if bestFitness <= s.Config.ThresholdValue {
					if newData, stop := s.Config.SuccessCallback(bestFitness, s.Config.MutableData); stop {
						return s.Population[0], newData
					} else {
						s.Config.MutableData = newData
					}
				}
			case sometinyai.Closest:
				if math.Abs(bestFitness-s.Config.ThresholdValue) < 0.0001 {
					if newData, stop := s.Config.SuccessCallback(bestFitness, s.Config.MutableData); stop {
						return s.Population[0], newData
					} else {
						s.Config.MutableData = newData
					}
				}
			}
		}

		fmt.Printf("Iteration %d | Fitness: %.4f | Data: %v\n",
			iter, bestFitness, s.Config.MutableData)
	}

	return s.Population[0], s.Config.MutableData
}
