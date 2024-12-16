package sometinyai

import (
	"fmt"
	"math/rand/v2"
	"os"

	"github.com/dominikbraun/graph"
)

type Genome struct {
	graph              graph.Graph[int, int]
	order              []int // This is the topological order of the nodes
	input              int
	output             int
	hidden             int
	adjacency          map[int]map[int]graph.Edge[int]
	activationFunction func(float64) float64 // This will be used for ALL nodes
}

type EdgeConnectionData struct {
	weight, bias float64
}

func NewGenome(x, y int, activation func(float64) float64) *Genome {
	g := graph.New(graph.IntHash, graph.Directed(), graph.Acyclic())
	for i := range x {
		g.AddVertex(i)
	}
	for i := range y {
		g.AddVertex(i + x)
	}
	for i := range x {
		for j := range y {
			g.AddEdge(i, j+x, graph.EdgeData(NewEdgeConnectionData(-1, -1)))
		}
	}
	return &Genome{
		graph:              g,
		order:              nil,
		input:              x,
		output:             y,
		hidden:             0,
		activationFunction: activation,
	}
}

func (g *Genome) ForwardPropagation(input ...float64) []float64 {
	// Check for correct input length
	if len(input) != g.input {
		panic(fmt.Sprintf("Expected %d inputs, got %d", g.input, len(input)))
	}

	// Get the topological ordering of the nodes

	if g.order == nil {
		g.order, _ = graph.TopologicalSort(g.graph)
	}
	ordering := g.order

	// Initialize node values
	nodeCount, _ := g.graph.Order()
	nodeValues := make([]float64, nodeCount)

	// Set input node values
	for i, v := range input {
		nodeValues[i] = v
	}
	// Iterate over nodes in topological order
	for _, node := range ordering {
		// Skip input nodes
		if node < g.input {
			continue
		}

		var sum float64

		// Get incoming edges to the current node
		if g.adjacency == nil {
			g.adjacency, _ = g.graph.AdjacencyMap()
		}
		adj := g.adjacency
		var inEdges []graph.Edge[int]
		for _, targets := range adj {
			if edge, exists := targets[node]; exists {
				inEdges = append(inEdges, edge)
			}
		}
		for _, edge := range inEdges {
			source := edge.Source
			edgeData := edge.Properties.Data.(*EdgeConnectionData)
			sum += nodeValues[source]*edgeData.weight + edgeData.bias
		}

		// Apply activation function (e.g., tanh)
		nodeValues[node] = g.activationFunction(sum)
	}

	// Collect output values
	outputValues := make([]float64, g.output)
	for i := 0; i < g.output; i++ {
		outputIndex := g.input + i
		outputValues[i] = nodeValues[outputIndex]
	}

	return outputValues
}

func NewEdgeConnectionData(weight, bias float64) *EdgeConnectionData {
	if weight < 0 {
		weight = rand.NormFloat64()
	}
	if bias < 0 {
		bias = rand.NormFloat64()
	}
	return &EdgeConnectionData{
		weight: weight,
		bias:   bias,
	}
}

func (g *Genome) Print() {
	adj, _ := g.graph.AdjacencyMap()
	for k, v := range adj {
		for k2, v2 := range v {
			fmt.Printf("%d -> %d: %v\n", k, k2, v2.Properties.Data.(*EdgeConnectionData))
		}
	}
}

func (g *Genome) Copy() *Genome {
	adj, _ := g.graph.AdjacencyMap()
	newGraph := graph.New(graph.IntHash, graph.Directed(), graph.Acyclic())

	// Copy vertices
	for source := range adj {
		newGraph.AddVertex(source)
	}

	// Copy edges and their data
	for source, targets := range adj {
		for target, edge := range targets {
			oldData := edge.Properties.Data.(*EdgeConnectionData)

			// Deep copy of EdgeConnectionData
			newData := &EdgeConnectionData{
				weight: oldData.weight,
				bias:   oldData.bias,
			}

			// Add the edge with the copied data
			newGraph.AddEdge(source, target, graph.EdgeData(newData))
		}
	}

	return &Genome{
		graph:              newGraph,
		input:              g.input,
		output:             g.output,
		hidden:             g.hidden,
		activationFunction: g.activationFunction,
	}
}

func (g *Genome) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	node_count, _ := g.graph.Order()
	fmt.Fprintf(file, "%d %d %d\n", node_count, g.input, g.output)
	adj, _ := g.graph.AdjacencyMap()
	for k, v := range adj {
		for k2, v2 := range v {
			data := v2.Properties.Data.(*EdgeConnectionData)
			fmt.Fprintf(file, "%d %d %f %f\n", k, k2, data.weight, data.bias)
		}
	}
	return nil
}
