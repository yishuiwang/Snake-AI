package search

import (
	"Snake-go/game"
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type Chromosome struct {
	// Genes is the weight of the neural network
	Genes []float64
	// score is the score of the snake
	Score float64
	// steps is the steps of the snake
	Steps float64
	// fitness is the fitness of the chromosome
	fitness float64
}

func NewChromosome() *Chromosome {
	var genes = make([]float64, 32*20+20*12+12*4)
	for i := range genes {
		genes[i] = rand.Float64()*2 - 1
	}
	return &Chromosome{Genes: genes}
}

func (c *Chromosome) Fitness(score int, steps int) float64 {
	c.Score, c.Steps = float64(score), float64(steps)
	c.fitness = (c.Score + 1/c.Steps) * 100000
	return c.fitness
}

// Mutate Gaussian mutation with scale of 0.2
func (c *Chromosome) Mutate() *Chromosome {
	var rate = 0.1
	for i := range c.Genes {
		if rand.Float64() < rate {
			c.Genes[i] = rand.Float64() - 0.2
		}
	}
	return c
}

// Crossover Single point crossover
func (c *Chromosome) Crossover(other *Chromosome) (*Chromosome, *Chromosome) {
	var pivot = rand.Intn(len(c.Genes))
	var offspring1 = c
	var offspring2 = other
	for i := pivot; i < len(c.Genes); i++ {
		offspring1.Genes[i], offspring2.Genes[i] = offspring2.Genes[i], offspring1.Genes[i]
	}
	return offspring1, offspring2
}

// ElitismSelect top size of the population
func ElitismSelect(population []*Chromosome) []*Chromosome {
	size := len(population) / 10
	SortByFitness(population)
	return population[:size]
}

// WheelSelect select the chromosome with the probability of fitness
func WheelSelect(population []*Chromosome, n int) []*Chromosome {
	var totalFitness float64
	for _, c := range population {
		totalFitness += c.fitness
	}
	newPopulation := make([]*Chromosome, n)
	for i := 0; i < n; i++ {
		r := rand.Float64() * totalFitness
		for _, c := range population {
			r -= c.fitness
			if r <= 0 {
				newPopulation[i] = c
				break
			}
		}
	}
	return newPopulation
}

func Evolve(population []*Chromosome) []*Chromosome {
	n := len(population)
	newPopulation := make([]*Chromosome, 0)
	// Elitism
	elites := ElitismSelect(population)
	newPopulation = append(newPopulation, elites...)
	// Wheel selection
	for i := len(elites); i < n; i++ {
		parents := WheelSelect(population, 2)
		offspring1, _ := parents[0].Crossover(parents[1])
		offspring1.Mutate()
		newPopulation = append(newPopulation, offspring1)
	}

	return newPopulation
}

// SortByFitness sort the population by fitness in descending order
func SortByFitness(population []*Chromosome) {
	sort.Slice(population, func(i, j int) bool {
		return population[i].fitness > population[j].fitness
	})
}

func Train(generation int, scale int) *Chromosome {
	games := make([]*game.PlayGround, scale)
	for i := 0; i < scale; i++ {
		games[i] = game.NewPlayGround()
	}
	nn := make([]*NeuralNetwork, scale)
	for i := 0; i < scale; i++ {
		nn[i] = NewNeuralNetwork()
	}
	chromosomes := make([]*Chromosome, scale)
	for i := 0; i < scale; i++ {
		chromosomes[i] = NewChromosome()
		nn[i].UpdateWeights(chromosomes[i].Genes)
	}

	var wg sync.WaitGroup
	var best *Chromosome
	// Training
	for g := 0; g < generation; g++ {
		// play game until all snakes are dead
		for i := 0; i < scale; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				for {
					state := games[i].GetState()
					dir := nn[i].Predict(state)
					if !games[i].Move(dir) {
						break
					}
				}
				score, steps := games[i].Score, games[i].Steps
				chromosomes[i].Fitness(score, steps)
			}(i)
		}
		wg.Wait()
		// evolve
		chromosomes = Evolve(chromosomes)
		best = chromosomes[0]
		for i := 0; i < scale; i++ {
			// update weights
			nn[i].UpdateWeights(chromosomes[i].Genes)
			// reset game
			games[i] = game.NewPlayGround()
		}
	}
	// save best to file
	return best
}

func SaveBest(best *Chromosome) {
	var file, _ = os.Create("./genes/best")
	defer file.Close()
	var writer = bufio.NewWriter(file)
	for _, w := range best.Genes {
		writer.WriteString(fmt.Sprintf("%.16f ", w))
	}
}

func LoadBest() *Chromosome {
	var file, _ = os.Open("./genes/best")
	defer file.Close()
	var reader = bufio.NewReader(file)
	var data, _ = reader.ReadString('\n')
	var genes = strings.Fields(data)
	var best = NewChromosome()
	for i := range genes {
		best.Genes[i], _ = strconv.ParseFloat(genes[i], 64)
	}
	return best
}
