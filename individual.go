package main

import (
	"math/rand"
)

type individual struct {
	dna []point
}

//mutate swaps the points around randomly, a hihger mutateProbability means more swaps
func mutate(dna []point, mutateProbability float32) {
	for i := 0; i < len(dna); i++ {
		if rand.Float32() < mutateProbability {
			i1 := rand.Intn(len(dna))
			dna[i], dna[i1] = dna[i1], dna[i]
		}
	}
}

func newIndividual(old *individual, mutateProbability float32) *individual {
	n := &individual{
		dna: make([]point, len(old.dna)),
	}
	copy(n.dna, old.dna)
	mutate(n.dna, mutateProbability)
	return n
}
