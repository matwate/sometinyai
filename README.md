# SomeTinyAI

A lightweight neural network framework using genetic algorithms for evolution based training

## Overview

SomeTinyAI provides a minimalist approach to neural networks, focusing on a graph-based genome representation that evoulves through mutations rather than traditional backpropagation

## Installation

```bash
go get github.com/matwate/sometinyai@latest
```

## Core Concepts

- **Genome**: Neural network represented as a directed graph
- **Mutation**: Evolves networks through adding/splitting connections and changing weight/bias values
- **Simulation**: Manages populations of genomes with fitness based seletion

# Usage

## Basic example

```go
package main

import (
    "fmt"

    "github.com/matwate/sometinyai"
    "github.com/matwate/sometinyai/activation"
    "github.com/matwate/sometinyai/simulation"
)

func main() {
    // Create a simulation with 2 inputs, 1 output
    sim := simulation.NewSimulation(
        2, // inputs
        1, // outputs
        activation.Relu, // activation function
        simulation.PopulationSize(100),
        simulation.Iterations(500),
        simulation.Fitness(myFitnessFunction),
    )

    // Train the network
    best, _ := sim.Train()

    // Use the trained network
    result := best.Genome.ForwardPropagation(1.0, 2.0)
    fmt.Println("Output:", result)
}

func myFitnessFunction(g *sometinyai.Genome, data interface{}) float64 {
    // Example: XOR fitness function
    inputs := [][]float64{
        {0, 0},
        {0, 1},
        {1, 0},
        {1, 1},
    }
    targets := []float64{0, 1, 1, 0}

    var fitness float64
    for i, input := range inputs {
        output := g.ForwardPropagation(input...)
        // Higher fitness = lower error
        fitness += 1.0 - abs(output[0] - targets[i])
    }
    return fitness
}

```

## Advanced Features

```go

// Set timeout for training
simulation.WithTimeout(5*time.Second)

// Use custom threshold for breeding selection
simulation.Threshold(simulation.Highest, 0.95)

// Save a trained network
genome.Save("mynetwork.genome", activation.Relu)

// Load a trained network
loadedGenome, _ := sometinyai.LoadGenome("mynetwork.genome")
```

```

```
