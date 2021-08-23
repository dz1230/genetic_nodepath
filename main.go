package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type point struct {
	x, y float32
}

func distance(route []point) float32 {
	var d float32
	l := len(route)
	for i := 1; i < l+1; i++ {
		prev := route[i-1]
		cur := route[i%l]
		xdist := cur.x - prev.x
		ydist := cur.y - prev.y
		d += float32(math.Sqrt(float64(xdist*xdist + ydist*ydist)))
	}
	return d
}

//randomPoint with coordinates between 0 and max.
func randomPoint(max float32) point {
	return point{
		x: (rand.Float32() - 0.5) * max,
		y: (rand.Float32() - 0.5) * max,
	}
}

func hasPoint(route []point, p point) bool {
	for i := 0; i < len(route); i++ {
		if p == route[i] {
			return true
		}
	}
	return false
}

func sameDNA(dna1, dna2 []point) bool {
	if len(dna1) != len(dna2) {
		return false
	}
	for i := 0; i < len(dna1); i++ {
		if dna1[i] != dna2[i] {
			return false
		}
	}
	return true
}

func randomRoute(n int, coordSystemSize float32) []point {
	r := make([]point, n)
	for i := 0; i < n; i++ {
		r[i] = randomPoint(coordSystemSize)
	}
	return r
}

func solve(points []point, generations int, populationSize int, mutationProbability, crossoverProbability float32) *individual {
	p := makePopulation(populationSize, crossoverProbability, mutationProbability, points)
	var fittest *individual
	var maxFitness float32
	for g := 0; g < generations; g++ {
		fitness := 1 / distance(p.fittest.dna)
		if fitness > maxFitness {
			fittest = p.fittest
			maxFitness = fitness
		}
		nextGen(p)
	}
	return fittest
}

func main() {
	fmt.Println("genetic algorithm: shortest route")
	fmt.Println("---------------------------------")
	rand.Seed(time.Now().UnixNano())
	pC := float32(0.3)
	pM := float32(0.005)
	populationSize := 100
	generations := 10000
	nodes := make([]point, 0, 20)
	if len(os.Args) == 1 || os.Args[1] == "--help" || os.Args[1] == "?" || os.Args[1] == "--h" || os.Args[1] == "h" || os.Args[1] == "help" || os.Args[1] == "--?" {
		fmt.Print("options:\n--pc [float] specifies crossover probability (0.0-1.0)\n--pm [float] specifies mutation probability (0.0-1.0)\n--popsize [uint] specifies population size (0-MaxInt)\n--ngens [uint] specifies number of generations (0-MaxInt)\n--nodes [[float],[float];[float],[float];...etc] specifies the nodes between which to calculate the path, as a semicolon-separated list of 2D points\n")
		return
	}
	args := os.Args[1:]
	for i := 0; i < len(args)-1; i += 2 {
		switch args[i] {
		case "--pc":
			parsedPC, err := strconv.ParseFloat(args[i+1], 32)
			if err != nil {
				panic(err)
			}
			pC = float32(parsedPC)
		case "--pm":
			parsedPM, err := strconv.ParseFloat(args[i+1], 32)
			if err != nil {
				panic(err)
			}
			pM = float32(parsedPM)
		case "--popsize":
			parsedPopSize, err := strconv.ParseInt(args[i+1], 10, 32)
			if err != nil {
				panic(err)
			}
			populationSize = int(parsedPopSize)
		case "--ngens":
			parsedNGens, err := strconv.ParseInt(args[i+1], 10, 32)
			if err != nil {
				panic(err)
			}
			generations = int(parsedNGens)
		case "--nodes":
			points := strings.Split(args[i+1], ";")
			for i := 0; i < len(points); i++ {
				coords := strings.Split(points[i], ",")
				parsedX, err := strconv.ParseFloat(coords[0], 32)
				if err != nil {
					panic(err)
				}
				parsedY, err := strconv.ParseFloat(coords[1], 32)
				if err != nil {
					panic(err)
				}
				nodes = append(nodes, point{x: float32(parsedX), y: float32(parsedY)})
			}
		}
	}
	if len(nodes) == 0 {
		nodes = randomRoute(cap(nodes), 50)
	}
	fmt.Printf("Population size: %d, generations: %d, crossover probability: %f, mutation probability: %f, nodes: %v\n", populationSize, generations, pC, pM, nodes)
	fmt.Println("Solving...")
	alltimeFittest := solve(nodes, generations, populationSize, pM, pC)
	fmt.Printf("Fittest: {distance: %f, fitness: %f, route: %v}\n", distance(alltimeFittest.dna), 1/distance(alltimeFittest.dna), alltimeFittest.dna)
}
