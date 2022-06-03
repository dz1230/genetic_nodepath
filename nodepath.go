package main

import (
	"errors"
	"math"
	"math/rand"
)

type point struct {
	x, y float32
}

type population struct {
	nodes                []point
	current              [][]int
	currentFitnesses     []float32
	fittest              []int
	dnaLength            int
	crossoverProbability float32
	mutationProbability  float32
	generation           int
}

func makePopulation(size int, crossoverProbability, mutationProbability float32, seed []point) *population {
	initialPop := make([][]int, size)
	for i := 0; i < size; i++ {
		p0individual := make([]int, len(seed))
		for j := 0; j < len(seed); j++ {
			p0individual[j] = j
		}
		shuffle(p0individual)
		initialPop[i] = p0individual
	}
	p := &population{
		nodes:                seed,
		current:              initialPop,
		currentFitnesses:     make([]float32, size),
		fittest:              make([]int, len(seed)),
		dnaLength:            len(seed),
		crossoverProbability: crossoverProbability,
		mutationProbability:  mutationProbability,
		generation:           0,
	}
	updateFitness(p)
	return p
}

func nextGen(p *population) {
	evolved := make([][]int, 0, len(p.current))
	for len(evolved) < len(p.current) {
		child1, child2 := twoChildren(p)
		evolved = append(evolved, child1, child2)
	}
	p.current = evolved
	updateFitness(p)
	p.generation++
}

//updateFitness fills p.currentFitnesses with current values and updates p.fittest
func updateFitness(p *population) {
	bestFitness, indexOfFittest := float32(0), -1
	for i := 0; i < len(p.current); i++ {
		f := 1 / distance(p.nodes, p.current[i])
		p.currentFitnesses[i] = f
		if i == 0 || bestFitness < f {
			indexOfFittest = i
			bestFitness = f
		}
	}
	if indexOfFittest >= 0 {
		copy(p.fittest, p.current[indexOfFittest])
	}
}

//twoChildren with parents from the current generation. The children are new slices.
func twoChildren(p *population) (child1, child2 []int) {
	mom, dad := roulette(p), roulette(p)
	if rand.Float32() < p.crossoverProbability {
		child1, child2 = crossover(p.current[mom], p.current[dad])
	} else {
		child1, child2 = make([]int, p.dnaLength), make([]int, p.dnaLength)
		copy(child1, p.current[mom])
		copy(child2, p.current[dad])
	}
	mutate(child1, p.mutationProbability)
	mutate(child2, p.mutationProbability)
	return
}

//roulette selection, preferring fit individuals over less fit individuals
func roulette(p *population) int {
	//take sum of fitnesses of all individuals
	fitnessSum := float32(0)
	for i := 0; i < len(p.currentFitnesses); i++ {
		fitnessSum += p.currentFitnesses[i]
	}
	//get random treshold
	roll := rand.Float32() * fitnessSum
	//select individual
	for i := 0; i < len(p.current); i++ {
		if roll < p.currentFitnesses[i] {
			return i
		}
		//decrease treshold so there will always be an individual selected
		roll -= p.currentFitnesses[i]
	}
	//should be practically impossible to reach
	panic(errors.New("Roulette selection did not select any individual"))
}

//crossover combines two routes by taking one (random) segment from one parent and the rest from the other parent. The children are new slices.
func crossover(mom, dad []int) (child1, child2 []int) {
	l := len(mom)
	//determine which part is copied from which parent
	num1 := rand.Intn(l)
	num2 := rand.Intn(l)
	segmentStart := num1
	segmentEnd := num2
	if num1 > num2 {
		segmentStart = num2
		segmentEnd = num1
	}
	//create offspring
	child1 = make([]int, l)
	copy(child1[segmentStart:segmentEnd], mom[segmentStart:segmentEnd])
	child2 = make([]int, l)
	copy(child2[segmentStart:segmentEnd], dad[segmentStart:segmentEnd])
	//copy points which are outside of the segment (while maintaining their order)
	j1, j2 := segmentEnd, segmentEnd
	for i := segmentEnd; i < l+segmentEnd; i++ {
		p := dad[i%l]
		if !hasPoint(child1[segmentStart:segmentEnd], p) {
			child1[j1%l] = p
			j1++
		}
		p = mom[i%l]
		if !hasPoint(child2[segmentStart:segmentEnd], p) {
			child2[j2%l] = p
			j2++
		}
	}
	return
}

//shuffle the order of points using fisher swap
func shuffle(path []int) {
	for i := len(path); i != 0; i-- {
		j := rand.Intn(i + 1)
		path[i], path[j] = path[j], path[i]
	}
}

//mutate swaps the points around randomly, a hihger mutationProbability means more swaps
func mutate(path []int, mutationProbability float32) {
	for i := 0; i < len(path); i++ {
		if rand.Float32() < mutationProbability {
			j := rand.Intn(len(path))
			path[i], path[j] = path[j], path[i]
		}
	}
}

//distance along the full path (including from last to first node)
func distance(nodes []point, path []int) float32 {
	var d float32
	l := len(path)
	for i := 1; i < l+1; i++ {
		prev := nodes[path[i-1]]
		cur := nodes[path[i%l]]
		xdist := cur.x - prev.x
		ydist := cur.y - prev.y
		d += float32(math.Sqrt(float64(xdist*xdist + ydist*ydist)))
	}
	return d
}

//hasPoint is true if path contains p
func hasPoint(path []int, p int) bool {
	for i := 0; i < len(path); i++ {
		if p == path[i] {
			return true
		}
	}
	return false
}

//randomPoint with coordinates between 0 and max.
func randomPoint(max float32) point {
	return point{
		x: (rand.Float32() - 0.5) * max,
		y: (rand.Float32() - 0.5) * max,
	}
}

//randomRoute with n random points
func randomRoute(n int, coordSystemSize float32) []point {
	r := make([]point, n)
	for i := 0; i < n; i++ {
		r[i] = randomPoint(coordSystemSize)
	}
	return r
}

//samePath is true if both paths are exactly the same
func samePath(path1, path2 []int) bool {
	if len(path1) != len(path2) {
		return false
	}
	for i := 0; i < len(path1); i++ {
		if path1[i] != path2[i] {
			return false
		}
	}
	return true
}
