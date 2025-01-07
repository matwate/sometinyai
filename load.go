package sometinyai

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/dominikbraun/graph"
)

func LoadGenome(filename string) *Genome {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	s, _ := io.ReadAll(file)
	lines := strings.Split(string(s), "\n")
	count_line := strings.Split(lines[0], " ")
	node_count, _ := strconv.Atoi(count_line[0])
	input, _ := strconv.Atoi(count_line[1])
	output, _ := strconv.Atoi(count_line[2])

	g := graph.New(graph.IntHash, graph.Directed(), graph.Acyclic())
	for i := range node_count {
		g.AddVertex(i)
	}
	for i := 1; i < len(lines); i++ {
		if len(lines[i]) == 0 {
			continue
		}
		line := strings.Split(lines[i], " ")
		fmt.Println(line)
		source, _ := strconv.Atoi(line[0])
		target, _ := strconv.Atoi(line[1])
		weight, _ := strconv.ParseFloat(line[2], 64)
		bias, _ := strconv.ParseFloat(line[3], 64)
		g.AddEdge(source, target, graph.EdgeData(&EdgeConnectionData{weight: weight, bias: bias}))
	}
	return &Genome{
		graph:  g,
		input:  input,
		output: output,
	}
}
