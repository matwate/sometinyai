package simulation

import (
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
	for iter := 0; iter < s.Config.iterations; iter++ {"type":"excalidraw/clipboard","elements":[{"id":"8Be7BY2xD9MsXc3XfchIY","type":"rectangle","x":1864.4126984126983,"y":-459.3809523809522,"width":827.142857142857,"height":650,"angle":0,"strokeColor":"#1e1e1e","backgroundColor":"transparent","fillStyle":"solid","strokeWidth":2,"strokeStyle":"solid","roughness":1,"opacity":100,"groupIds":[],"frameId":null,"index":"b1VO","roundness":{"type":3},"seed":1033814303,"version":715,"versionNonce":2108841215,"isDeleted":false,"boundElements":[],"updated":1738177923844,"link":null,"locked":false},{"id":"4OtBm9TaVr64MCcQhKRGi","type":"line","x":2301.5555555555557,"y":-463.6666666666664,"width":1.4285714285720132,"height":651.4285714285716,"angle":0,"strokeColor":"#1e1e1e","backgroundColor":"transparent","fillStyle":"solid","strokeWidth":2,"strokeStyle":"solid","roughness":1,"opacity":100,"groupIds":[],"frameId":null,"index":"b1Yl","roundness":{"type":2},"seed":25207889,"version":575,"versionNonce":645686559,"isDeleted":false,"boundElements":[],"updated":1738177923844,"link":null,"locked":false,"points":[[0,0],[1.4285714285720132,651.4285714285716]],"lastCommittedPoint":null,"startBinding":null,"endBinding":null,"startArrowhead":null,"endArrowhead":null},{"id":"esGnrM0vs2yZGfvkbJZfI","type":"line","x":1862.9841269841268,"y":-337.95238095238074,"width":828.5714285714287,"height":2.8571428571428896,"angle":0,"strokeColor":"#1e1e1e","backgroundColor":"transparent","fillStyle":"solid","strokeWidth":2,"strokeStyle":"solid","roughness":1,"opacity":100,"groupIds":[],"frameId":null,"index":"b1bl","roundness":{"type":2},"seed":1334265151,"version":576,"versionNonce":353397055,"isDeleted":false,"boundElements":[],"updated":1738177923844,"link":null,"locked":false,"points":[[0,0],[828.5714285714287,-2.8571428571428896]],"lastCommittedPoint":null,"startBinding":null,"endBinding":null,"startArrowhead":null,"endArrowhead":null},{"id":"DbmlyW_KaKxnVR5RWh3cQ","type":"text","x":1871.5555555555557,"y":-569.3809523809523,"width":490.4333190917969,"height":168.5714285714285,"angle":0,"strokeColor":"#1e1e1e","backgroundColor":"transparent","fillStyle":"solid","strokeWidth":2,"strokeStyle":"solid","roughness":1,"opacity":100,"groupIds":[],"frameId":null,"index":"b1cl","roundness":null,"seed":1131016753,"version":715,"versionNonce":1993271647,"isDeleted":false,"boundElements":[],"updated":1738177923844,"link":null,"locked":false,"text":"Tabla Servicios\n","fontSize":67.4285714285714,"fontFamily":5,"textAlign":"left","verticalAlign":"top","containerId":null,"originalText":"Tabla Servicios\n","autoResize":true,"lineHeight":1.25},{"id":"igKgBjYCA2VndzDeu6zqd","type":"text","x":1898.6984126984125,"y":-443.6666666666664,"width":338.6298941798943,"height":112.14285714285717,"angle":0,"strokeColor":"#1e1e1e","backgroundColor":"transparent","fillStyle":"solid","strokeWidth":2,"strokeStyle":"solid","roughness":1,"opacity":100,"groupIds":[],"frameId":null,"index":"b1dl","roundness":null,"seed":1980384607,"version":512,"versionNonce":2030593407,"isDeleted":false,"boundElements":[],"updated":1738177923844,"link":null,"locked":false,"text":"Columna","fontSize":89.7142857142857,"fontFamily":5,"textAlign":"left","verticalAlign":"top","containerId":null,"originalText":"Columna","autoResize":true,"lineHeight":1.25},{"id":"gG91DfmXNC0xwOsusr9Bx","type":"text","x":2457.269841269842,"y":-463.6666666666664,"width":36.983333587646484,"height":112.14285714285711,"angle":0,"strokeColor":"#1e1e1e","backgroundColor":"transparent","fillStyle":"solid","strokeWidth":2,"strokeStyle":"solid","roughness":1,"opacity":100,"groupIds":[],"frameId":null,"index":"b1fl","roundness":null,"seed":2058257425,"version":676,"versionNonce":1258118559,"isDeleted":false,"boundElements":[],"updated":1738177923844,"link":null,"locked":false,"text":"-","fontSize":89.7142857142857,"fontFamily":5,"textAlign":"left","verticalAlign":"top","containerId":null,"originalText":"-","autoResize":true,"lineHeight":1.25},{"id":"VBkLtMyrZp-u8fFyjiPn2","type":"line","x":1867.269841269841,"y":-202.23809523809518,"width":831.4285714285712,"height":1.4285714285713311,"angle":0,"strokeColor":"#1e1e1e","backgroundColor":"transparent","fillStyle":"solid","strokeWidth":2,"strokeStyle":"solid","roughness":1,"opacity":100,"groupIds":[],"frameId":null,"index":"b1hl","roundness":{"type":2},"seed":813187455,"version":599,"versionNonce":1608405439,"isDeleted":false,"boundElements":[],"updated":1738177923844,"link":null,"locked":false,"points":[[0,0],[831.4285714285712,1.4285714285713311]],"lastCommittedPoint":null,"startBinding":null,"endBinding":null,"startArrowhead":null,"endArrowhead":null},{"id":"j5d19CtY0bcFhSRpcHC2t","type":"line","x":1862.9841269841272,"y":-77.95238095238085,"width":832.8571428571428,"height":2.8571428571428896,"angle":0,"strokeColor":"#1e1e1e","backgroundColor":"transparent","fillStyle":"solid","strokeWidth":2,"strokeStyle":"solid","roughness":1,"opacity":100,"groupIds":[],"frameId":null,"index":"b1iC","roundness":{"type":2},"seed":1665522161,"version":704,"versionNonce":1892401631,"isDeleted":false,"boundElements":[],"updated":1738177923844,"link":null,"locked":false,"points":[[0,0],[832.8571428571428,-2.8571428571428896]],"lastCommittedPoint":null,"startBinding":null,"endBinding":null,"startArrowhead":null,"endArrowhead":null},{"id":"Qesdk-vcOVZVcI8n7s_oL","type":"line","x":1860.126984126984,"y":182.04761904761938,"width":827.142857142857,"height":10,"angle":0,"strokeColor":"#1e1e1e","backgroundColor":"transparent","fillStyle":"solid","strokeWidth":2,"strokeStyle":"solid","roughness":1,"opacity":100,"groupIds":[],"frameId":null,"index":"b1iQ","roundness":{"type":2},"seed":1074903455,"version":813,"versionNonce":1054223871,"isDeleted":false,"boundElements":[],"updated":1738177923844,"link":null,"locked":false,"points":[[0,0],[827.142857142857,10]],"lastCommittedPoint":null,"startBinding":null,"endBinding":null,"startArrowhead":null,"endArrowhead":null},{"id":"vjyNQIHTaSuFIPpaV5IB8","type":"line","x":1870.126984126984,"y":43.47619047619048,"width":832.8571428571428,"height":4.285714285714221,"angle":0,"strokeColor":"#1e1e1e","backgroundColor":"transparent","fillStyle":"solid","strokeWidth":2,"strokeStyle":"solid","roughness":1,"opacity":100,"groupIds":[],"frameId":null,"index":"b1it","roundness":{"type":2},"seed":1778855359,"version":752,"versionNonce":352196127,"isDeleted":false,"boundElements":[],"updated":1738177923845,"link":null,"locked":false,"points":[[0,0],[832.8571428571428,4.285714285714221]],"lastCommittedPoint":null,"startBinding":null,"endBinding":null,"startArrowhead":null,"endArrowhead":null},{"id":"SN32ClAKhQ2AYQxfAs4QC","type":"text","x":1891.5555555555557,"y":-309.38095238095195,"width":375.05,"height":67.81843026865872,"angle":0,"strokeColor":"#1e1e1e","backgroundColor":"transparent","fillStyle":"solid","strokeWidth":2,"strokeStyle":"solid","roughness":1,"opacity":100,"groupIds":[],"frameId":null,"index":"b1sl","roundness":null,"seed":2140852657,"version":722,"versionNonce":171178559,"isDeleted":false,"boundElements":[],"updated":1738177923845,"link":null,"locked":false,"text":"ID de servicio","fontSize":54.25474421492697,"fontFamily":5,"textAlign":"left","verticalAlign":"top","containerId":null,"originalText":"ID de servicio","autoResize":true,"lineHeight":1.25},{"id":"XMQ31KaecrFLEkPXIdwIy","type":"text","x":1894.4126984126983,"y":-173.66666666666617,"width":184.9,"height":67.81843026865872,"angle":0,"strokeColor":"#1e1e1e","backgroundColor":"transparent","fillStyle":"solid","strokeWidth":2,"strokeStyle":"solid","roughness":1,"opacity":100,"groupIds":[],"frameId":null,"index":"b1tl","roundness":null,"seed":2037980639,"version":826,"versionNonce":56121951,"isDeleted":false,"boundElements":[],"updated":1738177923845,"link":null,"locked":false,"text":"Nombre","fontSize":54.25474421492697,"fontFamily":5,"textAlign":"left","verticalAlign":"top","containerId":null,"originalText":"Nombre","autoResize":true,"lineHeight":1.25},{"id":"Ch_3fpPkylpF4hae4IAox","type":"text","x":1898.6984126984125,"y":-52.23809523809484,"width":297.1000061035156,"height":67.81843026865872,"angle":0,"strokeColor":"#1e1e1e","backgroundColor":"transparent","fillStyle":"solid","strokeWidth":2,"strokeStyle":"solid","roughness":1,"opacity":100,"groupIds":[],"frameId":null,"index":"b1ul","roundness":null,"seed":1284993937,"version":813,"versionNonce":1966847615,"isDeleted":false,"boundElements":[],"updated":1738177923845,"link":null,"locked":false,"text":"Descripcion","fontSize":54.25474421492697,"fontFamily":5,"textAlign":"left","verticalAlign":"top","containerId":null,"originalText":"Descripcion","autoResize":true,"lineHeight":1.25},{"id":"nZxJO2xGJjef79uIHrCaA","type":"text","x":1910.126984126984,"y":72.04761904761983,"width":164.43333333333334,"height":67.81843026865872,"angle":0,"strokeColor":"#1e1e1e","backgroundColor":"transparent","fillStyle":"solid","strokeWidth":2,"strokeStyle":"solid","roughness":1,"opacity":100,"groupIds":[],"frameId":null,"index":"b1vl","roundness":null,"seed":24644095,"version":885,"versionNonce":1721667231,"isDeleted":false,"boundElements":[],"updated":1738177923845,"link":null,"locked":false,"text":"Precio","fontSize":54.25474421492697,"fontFamily":5,"textAlign":"left","verticalAlign":"top","containerId":null,"originalText":"Precio","autoResize":true,"lineHeight":1.25},{"id":"Sh-PXUescpf9yyQCK21Fl","type":"text","x":2405.8412698412694,"y":-290.80952380952374,"width":172.66666666666666,"height":45,"angle":0,"strokeColor":"#1e1e1e","backgroundColor":"transparent","fillStyle":"solid","strokeWidth":2,"strokeStyle":"solid","roughness":1,"opacity":100,"groupIds":[],"frameId":null,"index":"b21V","roundness":null,"seed":1464501073,"version":454,"versionNonce":842097343,"isDeleted":false,"boundElements":[],"updated":1738177923845,"link":null,"locked":false,"text":"PRIMARY","fontSize":36,"fontFamily":5,"textAlign":"left","verticalAlign":"top","containerId":null,"originalText":"PRIMARY","autoResize":true,"lineHeight":1.25}],"files":{}}{
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
				return math.Abs(p[i].Fitness-s.Config.brea{"type":"excalidraw/clipboard","elements":[{"id":"8Be7BY2xD9MsXc3XfchIY","type":"rectangle","x":1864.4126984126983,"y":-459.3809523809522,"width":827.142857142857,"height":650,"angle":0,"strokeColor":"#1e1e1e","backgroundColor":"transparent","fillStyle":"solid","strokeWidth":2,"strokeStyle":"solid","roughness":1,"opacity":100,"groupIds":[],"frameId":null,"index":"b1VO","roundness":{"type":3},"seed":1033814303,"version":715,"versionNonce":2108841215,"isDeleted":false,"boundElements":[],"updated":1738177923844,"link":null,"locked":false},{"id":"4OtBm9TaVr64MCcQhKRGi","type":"line","x":2301.5555555555557,"y":-463.6666666666664,"width":1.4285714285720132,"height":651.4285714285716,"angle":0,"strokeColor":"#1e1e1e","backgroundColor":"transparent","fillStyle":"solid","strokeWidth":2,"strokeStyle":"solid","roughness":1,"opacity":100,"groupIds":[],"frameId":null,"index":"b1Yl","roundness":{"type":2},"seed":25207889,"version":575,"versionNonce":645686559,"isDeleted":false,"boundElements":[],"updated":1738177923844,"link":null,"locked":false,"points":[[0,0],[1.4285714285720132,651.4285714285716]],"lastCommittedPoint":null,"startBinding":null,"endBinding":null,"startArrowhead":null,"endArrowhead":null},{"id":"esGnrM0vs2yZGfvkbJZfI","type":"line","x":1862.9841269841268,"y":-337.95238095238074,"width":828.5714285714287,"height":2.8571428571428896,"angle":0,"strokeColor":"#1e1e1e","backgroundColor":"transparent","fillStyle":"solid","strokeWidth":2,"strokeStyle":"solid","roughness":1,"opacity":100,"groupIds":[],"frameId":null,"index":"b1bl","roundness":{"type":2},"seed":1334265151,"version":576,"versionNonce":353397055,"isDeleted":false,"boundElements":[],"updated":1738177923844,"link":null,"locked":false,"points":[[0,0],[828.5714285714287,-2.8571428571428896]],"lastCommittedPoint":null,"startBinding":null,"endBinding":null,"startArrowhead":null,"endArrowhead":null},{"id":"DbmlyW_KaKxnVR5RWh3cQ","type":"text","x":1871.5555555555557,"y":-569.3809523809523,"width":490.4333190917969,"height":168.5714285714285,"angle":0,"strokeColor":"#1e1e1e","backgroundColor":"transparent","fillStyle":"solid","strokeWidth":2,"strokeStyle":"solid","roughness":1,"opacity":100,"groupIds":[],"frameId":null,"index":"b1cl","roundness":null,"seed":1131016753,"version":715,"versionNonce":1993271647,"isDeleted":false,"boundElements":[],"updated":1738177923844,"link":null,"locked":false,"text":"Tabla Servicios\n","fontSize":67.4285714285714,"fontFamily":5,"textAlign":"left","verticalAlign":"top","containerId":null,"originalText":"Tabla Servicios\n","autoResize":true,"lineHeight":1.25},{"id":"igKgBjYCA2VndzDeu6zqd","type":"text","x":1898.6984126984125,"y":-443.6666666666664,"width":338.6298941798943,"height":112.14285714285717,"angle":0,"strokeColor":"#1e1e1e","backgroundColor":"transparent","fillStyle":"solid","strokeWidth":2,"strokeStyle":"solid","roughness":1,"opacity":100,"groupIds":[],"frameId":null,"index":"b1dl","roundness":null,"seed":1980384607,"version":512,"versionNonce":2030593407,"isDeleted":false,"boundElements":[],"updated":1738177923844,"link":null,"locked":false,"text":"Columna","fontSize":89.7142857142857,"fontFamily":5,"textAlign":"left","verticalAlign":"top","containerId":null,"originalText":"Columna","autoResize":true,"lineHeight":1.25},{"id":"gG91DfmXNC0xwOsusr9Bx","type":"text","x":2457.269841269842,"y":-463.6666666666664,"width":36.983333587646484,"height":112.14285714285711,"angle":0,"strokeColor":"#1e1e1e","backgroundColor":"transparent","fillStyle":"solid","strokeWidth":2,"strokeStyle":"solid","roughness":1,"opacity":100,"groupIds":[],"frameId":null,"index":"b1fl","roundness":null,"seed":2058257425,"version":676,"versionNonce":1258118559,"isDeleted":false,"boundElements":[],"updated":1738177923844,"link":null,"locked":false,"text":"-","fontSize":89.7142857142857,"fontFamily":5,"textAlign":"left","verticalAlign":"top","containerId":null,"originalText":"-","autoResize":true,"lineHeight":1.25},{"id":"VBkLtMyrZp-u8fFyjiPn2","type":"line","x":1867.269841269841,"y":-202.23809523809518,"width":831.4285714285712,"height":1.4285714285713311,"angle":0,"strokeColor":"#1e1e1e","backgroundColor":"transparent","fillStyle":"solid","strokeWidth":2,"strokeStyle":"solid","roughness":1,"opacity":100,"groupIds":[],"frameId":null,"index":"b1hl","roundness":{"type":2},"seed":813187455,"version":599,"versionNonce":1608405439,"isDeleted":false,"boundElements":[],"updated":1738177923844,"link":null,"locked":false,"points":[[0,0],[831.4285714285712,1.4285714285713311]],"lastCommittedPoint":null,"startBinding":null,"endBinding":null,"startArrowhead":null,"endArrowhead":null},{"id":"j5d19CtY0bcFhSRpcHC2t","type":"line","x":1862.9841269841272,"y":-77.95238095238085,"width":832.8571428571428,"height":2.8571428571428896,"angle":0,"strokeColor":"#1e1e1e","backgroundColor":"transparent","fillStyle":"solid","strokeWidth":2,"strokeStyle":"solid","roughness":1,"opacity":100,"groupIds":[],"frameId":null,"index":"b1iC","roundness":{"type":2},"seed":1665522161,"version":704,"versionNonce":1892401631,"isDeleted":false,"boundElements":[],"updated":1738177923844,"link":null,"locked":false,"points":[[0,0],[832.8571428571428,-2.8571428571428896]],"lastCommittedPoint":null,"startBinding":null,"endBinding":null,"startArrowhead":null,"endArrowhead":null},{"id":"Qesdk-vcOVZVcI8n7s_oL","type":"line","x":1860.126984126984,"y":182.04761904761938,"width":827.142857142857,"height":10,"angle":0,"strokeColor":"#1e1e1e","backgroundColor":"transparent","fillStyle":"solid","strokeWidth":2,"strokeStyle":"solid","roughness":1,"opacity":100,"groupIds":[],"frameId":null,"index":"b1iQ","roundness":{"type":2},"seed":1074903455,"version":813,"versionNonce":1054223871,"isDeleted":false,"boundElements":[],"updated":1738177923844,"link":null,"locked":false,"points":[[0,0],[827.142857142857,10]],"lastCommittedPoint":null,"startBinding":null,"endBinding":null,"startArrowhead":null,"endArrowhead":null},{"id":"vjyNQIHTaSuFIPpaV5IB8","type":"line","x":1870.126984126984,"y":43.47619047619048,"width":832.8571428571428,"height":4.285714285714221,"angle":0,"strokeColor":"#1e1e1e","backgroundColor":"transparent","fillStyle":"solid","strokeWidth":2,"strokeStyle":"solid","roughness":1,"opacity":100,"groupIds":[],"frameId":null,"index":"b1it","roundness":{"type":2},"seed":1778855359,"version":752,"versionNonce":352196127,"isDeleted":false,"boundElements":[],"updated":1738177923845,"link":null,"locked":false,"points":[[0,0],[832.8571428571428,4.285714285714221]],"lastCommittedPoint":null,"startBinding":null,"endBinding":null,"startArrowhead":null,"endArrowhead":null},{"id":"SN32ClAKhQ2AYQxfAs4QC","type":"text","x":1891.5555555555557,"y":-309.38095238095195,"width":375.05,"height":67.81843026865872,"angle":0,"strokeColor":"#1e1e1e","backgroundColor":"transparent","fillStyle":"solid","strokeWidth":2,"strokeStyle":"solid","roughness":1,"opacity":100,"groupIds":[],"frameId":null,"index":"b1sl","roundness":null,"seed":2140852657,"version":722,"versionNonce":171178559,"isDeleted":false,"boundElements":[],"updated":1738177923845,"link":null,"locked":false,"text":"ID de servicio","fontSize":54.25474421492697,"fontFamily":5,"textAlign":"left","verticalAlign":"top","containerId":null,"originalText":"ID de servicio","autoResize":true,"lineHeight":1.25},{"id":"XMQ31KaecrFLEkPXIdwIy","type":"text","x":1894.4126984126983,"y":-173.66666666666617,"width":184.9,"height":67.81843026865872,"angle":0,"strokeColor":"#1e1e1e","backgroundColor":"transparent","fillStyle":"solid","strokeWidth":2,"strokeStyle":"solid","roughness":1,"opacity":100,"groupIds":[],"frameId":null,"index":"b1tl","roundness":null,"seed":2037980639,"version":826,"versionNonce":56121951,"isDeleted":false,"boundElements":[],"updated":1738177923845,"link":null,"locked":false,"text":"Nombre","fontSize":54.25474421492697,"fontFamily":5,"textAlign":"left","verticalAlign":"top","containerId":null,"originalText":"Nombre","autoResize":true,"lineHeight":1.25},{"id":"Ch_3fpPkylpF4hae4IAox","type":"text","x":1898.6984126984125,"y":-52.23809523809484,"width":297.1000061035156,"height":67.81843026865872,"angle":0,"strokeColor":"#1e1e1e","backgroundColor":"transparent","fillStyle":"solid","strokeWidth":2,"strokeStyle":"solid","roughness":1,"opacity":100,"groupIds":[],"frameId":null,"index":"b1ul","roundness":null,"seed":1284993937,"version":813,"versionNonce":1966847615,"isDeleted":false,"boundElements":[],"updated":1738177923845,"link":null,"locked":false,"text":"Descripcion","fontSize":54.25474421492697,"fontFamily":5,"textAlign":"left","verticalAlign":"top","containerId":null,"originalText":"Descripcion","autoResize":true,"lineHeight":1.25},{"id":"nZxJO2xGJjef79uIHrCaA","type":"text","x":1910.126984126984,"y":72.04761904761983,"width":164.43333333333334,"height":67.81843026865872,"angle":0,"strokeColor":"#1e1e1e","backgroundColor":"transparent","fillStyle":"solid","strokeWidth":2,"strokeStyle":"solid","roughness":1,"opacity":100,"groupIds":[],"frameId":null,"index":"b1vl","roundness":null,"seed":24644095,"version":885,"versionNonce":1721667231,"isDeleted":false,"boundElements":[],"updated":1738177923845,"link":null,"locked":false,"text":"Precio","fontSize":54.25474421492697,"fontFamily":5,"textAlign":"left","verticalAlign":"top","containerId":null,"originalText":"Precio","autoResize":true,"lineHeight":1.25},{"id":"Sh-PXUescpf9yyQCK21Fl","type":"text","x":2405.8412698412694,"y":-290.80952380952374,"width":172.66666666666666,"height":45,"angle":0,"strokeColor":"#1e1e1e","backgroundColor":"transparent","fillStyle":"solid","strokeWidth":2,"strokeStyle":"solid","roughness":1,"opacity":100,"groupIds":[],"frameId":null,"index":"b21V","roundness":null,"seed":1464501073,"version":454,"versionNonce":842097343,"isDeleted":false,"boundElements":[],"updated":1738177923845,"link":null,"locked":false,"text":"PRIMARY","fontSize":36,"fontFamily":5,"textAlign":"left","verticalAlign":"top","containerId":null,"originalText":"PRIMARY","autoResize":true,"lineHeight":1.25}],"files":{}}kValue) < math.Abs(p[j].Fitness-s.Config.breakValue)
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
