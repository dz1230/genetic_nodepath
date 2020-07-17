package main

import (
	"math"
	"math/rand"
)

type population struct {
	current              []*individual
	currentFitnesses     []float32
	crossoverProbability float32
	mutationProbability  float32
	generation           int
	fittest              *individual
}

func makePopulation(size int, crossoverProbability, mutationProbability float32, seed []point) *population {
	initialPop := make([]*individual, size)
	for i := 0; i < size; i++ {
		p0individual := &individual{dna: make([]point, len(seed))}
		copy(p0individual.dna, seed)
		shuffle(p0individual.dna, fisherSwap)
		initialPop[i] = p0individual
	}
	p := &population{
		current:              initialPop,
		currentFitnesses:     make([]float32, size),
		crossoverProbability: crossoverProbability,
		mutationProbability:  mutationProbability,
		generation:           0,
		fittest:              nil,
	}
	for i := 0; i < len(p.current); i++ {
		p.currentFitnesses[i] = 1 / distance(p.current[i].dna)
	}
	p.fittest = getFittest(p)
	return p
}

func nextGen(p *population) {
	evolved := make([]*individual, 0, len(p.current))
	for len(evolved) < len(p.current) {
		child1, child2 := twoChildren(p)
		evolved = append(evolved, child1, child2)
	}
	p.current = evolved
	for i := 0; i < len(p.current); i++ {
		p.currentFitnesses[i] = 1 / distance(p.current[i].dna)
	}
	p.fittest = getFittest(p)
	p.generation++
}

func twoChildren(p *population) (child1, child2 *individual) {
	mom := roulette(p)
	dad := roulette(p)
	if rand.Float32() < p.crossoverProbability {
		mom, dad = crossover(mom, dad)
	}
	return newIndividual(mom, p.mutationProbability), newIndividual(dad, p.mutationProbability)
}

func roulette(p *population) *individual {
	fitnessSum := float32(0)
	for i := 0; i < len(p.currentFitnesses); i++ {
		fitnessSum += p.currentFitnesses[i]
	}
	roll := rand.Float32() * fitnessSum
	for i := 0; i < len(p.current); i++ {
		if roll < p.currentFitnesses[i] {
			return p.current[i]
		}
		roll -= p.currentFitnesses[i]
	}
	return nil
}

func crossover(mom, dad *individual) (firstOffspring, secondOffspring *individual) {
	num1 := rand.Intn(len(mom.dna))
	num2 := rand.Intn(len(dad.dna))
	segmentStart := num1
	segementEnd := num2
	if num1 > num2 {
		segmentStart = num2
		segementEnd = num1
	}
	firstOffspring = &individual{dna: orderedCrossover(segmentStart, segementEnd, mom, dad)}
	secondOffspring = &individual{dna: orderedCrossover(segmentStart, segementEnd, dad, mom)}
	return
}

func orderedCrossover(startIndex, endIndex int, segParent, otherParent *individual) []point {
	childDNA := make([]point, endIndex-startIndex)
	copy(childDNA, segParent.dna[startIndex:endIndex])
	for i := 0; i < len(otherParent.dna); i++ {
		parentIndex := (endIndex + i) % len(otherParent.dna)
		parentLoc := otherParent.dna[parentIndex]
		if !hasPoint(childDNA, parentLoc) {
			childDNA = append(childDNA, parentLoc)
		}
	}
	return rotate(childDNA, startIndex)
}

func rotate(arr []point, i int) []point {
	n := make([]point, len(arr))
	offset := len(arr) - i
	ni := 0
	for i := offset; i < len(arr); i++ {
		n[ni] = arr[i]
		ni++
	}
	for i := 0; i < offset; i++ {
		n[ni] = arr[i]
		ni++
	}
	return n
}

func shuffle(arr []point, swapper func(arr []point, index int)) {
	i := len(arr)
	for i != 0 {
		i--
		swapper(arr, i)
	}
}

func fisherSwap(arr []point, index int) {
	r := rand.Intn(index + 1)
	arr[index], arr[r] = arr[r], arr[index]
}

func getFittest(p *population) *individual {
	fittestI := -1
	highestFitness := float32(-math.MaxFloat32)
	for i := 0; i < len(p.currentFitnesses); i++ {
		if p.currentFitnesses[i] > highestFitness {
			fittestI = i
			highestFitness = p.currentFitnesses[i]
		}
	}
	if fittestI < 0 {
		return nil
	}
	return p.current[fittestI]
}
