package activation

import "math"

func Tanh(x float64) float64 {
	return math.Tanh(x)
}

func Sigmoid(x float64) float64 {
	return 1 / (1 + math.Exp(-x))
}

func Relu(x float64) float64 {
	if x < 0 {
		return 0
	}
	return x
}

func LeakyRelu(x float64) float64 {
	if x < 0 {
		return 0.01 * x
	}
	return x
}
