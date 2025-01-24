package sometinyai

import (
	"fmt"
	"os"

	"github.com/dominikbraun/graph"
	"google.golang.org/protobuf/proto"

	"github.com/matwate/sometinyai/activation"
	pb "github.com/matwate/sometinyai/protos"
)

func (g *Genome) Save(filename string, act activation.ActivationFunction) error {
	genome := &pb.Genome{}
	genome.Inputs = int32(g.input)
	genome.Outputs = int32(g.output)
	genome.Neurons = int32(g.input + g.output + g.hidden)
	connections := []*pb.Connection{}
	adj, _ := g.graph.AdjacencyMap()
	for source, targets := range adj {
		for target, edge := range targets {
			data := edge.Properties.Data.(*EdgeConnectionData)
			conn := &pb.Connection{
				In:     int32(source),
				Out:    int32(target),
				Weight: float64(data.weight),
				Bias:   float64(data.bias),
			}
			connections = append(connections, conn)
		}
	}
	genome.Connections = connections
	genome.Activation = map[activation.ActivationFunction]string{
		activation.Sigmoid_T:   "Sigmoid",
		activation.LeakyRelu_T: "LeakyRelu",
		activation.Relu_T:      "Relu",
		activation.Tanh_T:      "Tanh",
	}[act]
	out, err := proto.Marshal(genome)
	fmt.Printf("Size of the genome: %d Bytes\n", len(out))
	if err != nil {
		return err
	}
	if err := os.WriteFile(filename, out, 0644); err != nil {
		return err
	}
	return nil
}

func LoadGenome(filename string) (*Genome, error) {
	in, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	genome := &pb.Genome{}
	if err := proto.Unmarshal(in, genome); err != nil {
		return nil, err
	}
	return toGenome(genome), nil
}

func toGenome(genome *pb.Genome) *Genome {
	act := map[string]func(float64) float64{
		"Sigmoid":   activation.Sigmoid,
		"LeakyRelu": activation.LeakyRelu,
		"Relu":      activation.Relu,
		"Tanh":      activation.Tanh,
	}

	gr := graph.New(graph.IntHash, graph.Directed(), graph.Acyclic())

	for i := range genome.GetInputs() {
		gr.AddVertex(int(i))
	}
	for i := range genome.GetOutputs() {
		gr.AddVertex(int(i + genome.GetInputs()))
	}

	for i := range genome.GetNeurons() {
		gr.AddVertex(int(i + genome.GetInputs() + genome.GetOutputs()))
	}

	for _, conn := range genome.GetConnections() {
		gr.AddEdge(int(conn.GetIn()), int(conn.GetOut()), graph.EdgeData(&EdgeConnectionData{
			weight: conn.GetWeight(),
			bias:   conn.GetBias(),
		}))
	}

	g := &Genome{
		graph:              gr,
		input:              int(genome.GetInputs()),
		output:             int(genome.GetOutputs()),
		hidden:             int(genome.GetNeurons()),
		activationFunction: act[genome.GetActivation()],
		order:              nil,
	}
	return g
}
