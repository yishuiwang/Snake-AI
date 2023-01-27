package search

import (
	"math"
	"math/rand"
	"time"
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
	weight [][]float64
}

// NewNeuralNetwork 初始化神经网络结构 [32,20,12,4]
func NewNeuralNetwork() *NeuralNetwork {
	rand.Seed(time.Now().UnixNano())
	n := &NeuralNetwork{
		input:  make([]SigmoidNeuron, 32),
		hidden: make([][]SigmoidNeuron, 2),
		output: make([]SigmoidNeuron, 4),
	}
	n.hidden[0] = make([]SigmoidNeuron, 20)
	n.hidden[1] = make([]SigmoidNeuron, 12)
	n.weight = make([][]float64, 3)
	n.weight[0] = make([]float64, 32)
	n.weight[1] = make([]float64, 20)
	n.weight[2] = make([]float64, 12)
	for i := 0; i < 32; i++ {
		n.weight[0][i] = rand.Float64() - 0.5
	}
	for i := 0; i < 20; i++ {
		n.weight[1][i] = rand.Float64() - 0.5
	}
	for i := 0; i < 12; i++ {
		n.weight[2][i] = rand.Float64() - 0.5
	}
	return n
}

func (nn *NeuralNetwork) Predict(input []float64) int {
	// 输入层
	for i := 0; i < 32; i++ {
		nn.input[i].inputSignal = input[i]
		nn.input[i].outputSignal = input[i]
	}
	// 输入层到隐藏层
	for i := 0; i < 20; i++ {
		nn.hidden[0][i].inputSignal = 0
		for j := 0; j < 32; j++ {
			nn.hidden[0][i].inputSignal += nn.input[j].outputSignal * nn.weight[0][j]
		}
		nn.hidden[0][i].outputSignal = nn.hidden[0][i].ReLU(nn.hidden[0][i].inputSignal)
	}
	for i := 0; i < 12; i++ {
		nn.hidden[1][i].inputSignal = 0
		for j := 0; j < 20; j++ {
			nn.hidden[1][i].inputSignal += nn.hidden[0][j].outputSignal * nn.weight[1][j]
		}
		nn.hidden[1][i].outputSignal = nn.hidden[1][i].ReLU(nn.hidden[1][i].inputSignal)
	}
	// 隐层到输出层
	for i := 0; i < 4; i++ {
		nn.output[i].inputSignal = 0
		for j := 0; j < 12; j++ {
			nn.output[i].inputSignal += nn.hidden[1][j].outputSignal * nn.weight[2][j]
		}
		nn.output[i].outputSignal = nn.output[i].Sigmoid(nn.output[i].inputSignal)
	}
	// 找到最大输出
	max := 0
	for i := 0; i < 4; i++ {
		if nn.output[i].outputSignal > nn.output[max].outputSignal {
			max = i
		}
	}
	return max
}
