package search

import (
	"math/rand"
	"sort"
)

type Chromosome struct {
	// genes is the weight of the neural network
	genes []float64
	// score is the score of the snake
	score float64
	// steps is the steps of the snake
	steps int
	// fitness is the fitness of the chromosome
	fitness float64
}

func (c Chromosome) Fitness() float64 {
	c.score, c.steps = 0, 0
	c.fitness = (c.score + 1) / float64(c.steps) * 1000
	return c.fitness
}

// Mutate Gaussian mutation with scale of 0.2
func (c Chromosome) Mutate() Chromosome {
	for i := range c.genes {
		c.genes[i] += rand.NormFloat64() * 0.2
	}
	return c
}

// Crossover Single point crossover
func (c Chromosome) Crossover(other Chromosome) (Chromosome, Chromosome) {
	var pivot = rand.Intn(len(c.genes))
	var offspring1 = c
	var offspring2 = other
	for i := pivot; i < len(c.genes); i++ {
		offspring1.genes[i], offspring2.genes[i] = offspring2.genes[i], offspring1.genes[i]
	}
	return offspring1, offspring2
}

// ElitismSelect top size of the population
func ElitismSelect(population []Chromosome) []Chromosome {
	size := len(population) * 2 / 3
	sort.Slice(population, func(i, j int) bool {
		return population[i].fitness > population[j].fitness
	})
	return population[:size]
}

// WheelSelect select the chromosome with the probability of fitness
func WheelSelect(population []Chromosome, n int) []Chromosome {
	var totalFitness float64
	for _, c := range population {
		totalFitness += c.fitness
	}
	newPopulation := make([]Chromosome, n)
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

func Evolve(population []Chromosome) []Chromosome {
	n := len(population)
	// Elitism
	population = ElitismSelect(population)
	// Wheel Select && Crossover
	for len(population) < n {
		var parent1 = WheelSelect(population, 1)[0]
		var parent2 = WheelSelect(population, 1)[0]
		var offspring1, offspring2 = parent1.Crossover(parent2)
		offspring1 = offspring1.Mutate()
		offspring2 = offspring2.Mutate()
		population = append(population, offspring1, offspring2)
	}
	// reshuffle
	for i := range population {
		j := rand.Intn(i + 1)
		population[i], population[j] = population[j], population[i]
	}
	return population
}
