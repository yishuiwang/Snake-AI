package search

import (
	"math"
	"math/rand"
)

type SigmoidNeuron struct {
	inputSignal  float64
	outputSignal float64
}

func (n *SigmoidNeuron) Sigmoid(input float64) float64 {
	n.inputSignal = input
	n.outputSignal = 1.0 / (1.0 + math.Exp2(-input))
	return n.outputSignal
}

func (n *SigmoidNeuron) ReLU(input float64) float64 {
	n.inputSignal = input
	if input < 0 {
		n.outputSignal = 0
	} else {
		n.outputSignal = input
	}
	return n.outputSignal
}

type NeuralNetwork struct {
	input  []SigmoidNeuron
	hidden [][]SigmoidNeuron
	output []SigmoidNeuron
	w0     [][]float64
	w1     [][]float64
	w2     [][]float64
}

// NewNeuralNetwork 初始化神经网络结构 [32,20,12,4]
func NewNeuralNetwork() *NeuralNetwork {
	n := &NeuralNetwork{
		input:  make([]SigmoidNeuron, 32),
		hidden: make([][]SigmoidNeuron, 2),
		output: make([]SigmoidNeuron, 4),
	}
	n.hidden[0] = make([]SigmoidNeuron, 20)
	n.hidden[1] = make([]SigmoidNeuron, 12)
	n.w0 = make([][]float64, 20)
	n.w1 = make([][]float64, 12)
	n.w2 = make([][]float64, 4)
	for i := 0; i < 20; i++ {
		n.w0[i] = make([]float64, 32)
		for j := 0; j < 32; j++ {
			n.w0[i][j] = rand.Float64()*2 - 1
		}
	}
	for i := 0; i < 12; i++ {
		n.w1[i] = make([]float64, 20)
		for j := 0; j < 20; j++ {
			n.w1[i][j] = rand.Float64()*2 - 1
		}
	}
	for i := 0; i < 4; i++ {
		n.w2[i] = make([]float64, 12)
		for j := 0; j < 12; j++ {
			n.w2[i][j] = rand.Float64()*2 - 1
		}
	}
	return n
}

func (nn *NeuralNetwork) Predict(input []float64) int {
	// 输入层
	for i := 0; i < 32; i++ {
		nn.input[i].inputSignal = input[i]
		nn.input[i].outputSignal = input[i]
	}
	for i := 0; i < 32; i++ {
	}
	// 输入层到隐藏层
	for i := 0; i < 20; i++ {
		var sum float64
		for j := 0; j < 32; j++ {
			sum += nn.input[j].outputSignal * nn.w0[i][j]
		}
		nn.hidden[0][i].inputSignal = sum
		nn.hidden[0][i].outputSignal = nn.hidden[0][i].ReLU(sum)
	}
	// 隐藏层到隐藏层
	for i := 0; i < 12; i++ {
		var sum float64
		for j := 0; j < 20; j++ {
			sum += nn.hidden[0][j].outputSignal * nn.w1[i][j]
		}
		nn.hidden[1][i].inputSignal = sum
		nn.hidden[1][i].outputSignal = nn.hidden[1][i].ReLU(sum)
	}
	// 隐层到输出层
	for i := 0; i < 4; i++ {
		var sum float64
		for j := 0; j < 12; j++ {
			sum += nn.hidden[1][j].outputSignal * nn.w2[i][j]
		}
		nn.output[i].inputSignal = sum
		nn.output[i].outputSignal = nn.output[i].Sigmoid(sum)
	}
	// 输出层
	var max float64
	var index int
	for i := 0; i < 4; i++ {
		if nn.output[i].outputSignal > max {
			max = nn.output[i].outputSignal
			index = i
		}
	}
	return index
}

// UpdateWeights Update 更新神经网络权重
func (nn *NeuralNetwork) UpdateWeights(weights []float64) {
	index := 0
	for i := 0; i < 20; i++ {
		for j := 0; j < 32; j++ {
			nn.w0[i][j] = weights[index]
			index++
		}
	}
	for i := 0; i < 12; i++ {
		for j := 0; j < 20; j++ {
			nn.w1[i][j] = weights[index]
			index++
		}
	}
	for i := 0; i < 4; i++ {
		for j := 0; j < 12; j++ {
			nn.w2[i][j] = weights[index]
			index++
		}
	}
}
