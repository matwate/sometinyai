package sometinyai

import (
	"math/rand/v2"

	"github.com/dominikbraun/graph"
)

func (g *Genome) Mutate(count int) {
	// 0.1 chance to split a connection
	// 0.2 chance to add a connection
	// 0.5 chance to change a weight
	// 0.2 chance to change a bias
	// They all can happen at the same time
	for i := 0; i < count; i++ {
		n := rand.Float64()
		if n < 0.1 {
			g.SplitConnection()
		}
		if n < 0.2 {
			g.AddConnection()
		}
		if n < 0.5 {
			g.ChangeWeight()
		}
		if n < 0.2 {
			g.ChangeBias()
		}
	}
	g.order = nil // Cache invalidation???
}

func (g *Genome) SplitConnection() {
	edges, _ := g.graph.AdjacencyMap()
	// Find a non output node:
	n := rand.IntN(g.input + g.hidden)
	if n >= g.input {
		n += g.output
	}
	node := edges[n]
	if len(node) == 0 {
		return
	}
	edge := RandomValueOfMap(node)
	from, to := edge.Source, edge.Target
	weight, bias := edge.Properties.Data.(*EdgeConnectionData).weight, edge.Properties.Data.(*EdgeConnectionData).bias
	g.graph.AddVertex(g.input + g.output + g.hidden)
	g.hidden++
	g.graph.AddEdge(from, g.hidden+g.input+g.output-1, graph.EdgeData(NewEdgeConnectionData(1, 0)))
	g.graph.AddEdge(
		g.hidden+g.input+g.output-1,
		to,
		graph.EdgeData(NewEdgeConnectionData(weight, bias)),
	)
	g.graph.RemoveEdge(from, to)
}

func (g *Genome) AddConnection() {
	edges, _ := g.graph.AdjacencyMap()
	// Find a non output node:
	n := rand.IntN(g.input + g.hidden)
	if n >= g.input {
		n += g.output
	}
	node := edges[n]
	if len(node) == 0 {
		return
	}
	edge := RandomValueOfMap(node)
	from, to := edge.Source, edge.Target
	err := g.graph.AddEdge(from, to, graph.EdgeData(NewEdgeConnectionData(-1, -1)))
	if err != nil {
		return // we assume that the map is full
	}
}

func (g *Genome) ChangeWeight() {
	edges, _ := g.graph.AdjacencyMap()
	// Find a non output node:
	n := rand.IntN(g.input + g.hidden)
	if n >= g.input {
		n += g.output
	}
	node := edges[n]
	if len(node) == 0 {
		return
	}
	edge := RandomValueOfMap(node)
	edge.Properties.Data.(*EdgeConnectionData).weight = edge.Properties.Data.(*EdgeConnectionData).weight + rand.NormFloat64()
}

func (g *Genome) ChangeBias() {
	edges, _ := g.graph.AdjacencyMap()
	// Find a non output node:
	n := rand.IntN(g.input + g.hidden)
	if n >= g.input {
		n += g.output
	}
	node := edges[n]
	if len(node) == 0 {
		return
	}
	edge := RandomValueOfMap(node)

	edge.Properties.Data.(*EdgeConnectionData).bias = edge.Properties.Data.(*EdgeConnectionData).bias + rand.NormFloat64()
}

func RandomValueOfMap[T comparable, Y any](m map[T]Y) Y {
	if len(m) == 0 {
		panic("map is empty")
	}
	k := rand.IntN(len(m))
	for _, face := range m {
		if k == 0 {
			return face
		}
		k--
	}
	panic("unreachable")
}
