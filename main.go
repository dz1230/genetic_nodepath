package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func solve(nodes []point, generations int, populationSize int, mutationProbability, crossoverProbability float32) (shortestRoute []point) {
	p := makePopulation(populationSize, crossoverProbability, mutationProbability, nodes)
	var fittest []int
	var maxFitness float32
	for g := 0; g < generations; g++ {
		fitness := 1 / distance(p.nodes, p.fittest)
		if g == 0 || fitness > maxFitness {
			fittest = p.fittest
			maxFitness = fitness
		}
		nextGen(p)
	}
	shortestRoute = make([]point, len(nodes))
	for i := 0; i < len(fittest); i++ {
		shortestRoute[i] = nodes[fittest[i]]
	}
	return
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
	forwardPath := make([]int, len(alltimeFittest))
	for i := 0; i < len(alltimeFittest); i++ {
		forwardPath[i] = i
	}
	fmt.Printf("Fittest: {distance: %f, fitness: %f, route: %v}\n", distance(alltimeFittest, forwardPath), 1/distance(alltimeFittest, forwardPath), alltimeFittest)
}
