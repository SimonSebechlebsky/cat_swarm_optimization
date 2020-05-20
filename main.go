package main

import (
	"flag"
	"github.com/SimonSebechlebsky/cat_swarm_optimization/cso"
	"github.com/SimonSebechlebsky/cat_swarm_optimization/problem"
)

func main() {

	inputFile := flag.String("input_file", "./problem/inputs/b_should_be_easy.in", "Path to input file")
	catNum := flag.Int("cat_num", 20, "Number of cats in optimizer")
	mr := flag.Float64("mr", 0.7, "Mixture ratio between tracing mode and seeking mode")
	smp := flag.Int("smp", 50, "Seeking memory pool (how many copies in seeking mode)")
	srd := flag.Int("srd", 5, "Seeking range distance")
	velLimitflag := flag.Int("velocity_limit", 10, "Velocity limit")
	iterations := flag.Int("iterations", 1000, "Number of iterations")

	flag.Parse()

	problemDef := problem.LoadCarProblemDefinition(*inputFile)
	fitnessFunc := problem.GetFitnessFunc(problemDef)
	stateGenerator := problem.GetSolutionStateFunc(problemDef)

	optimizer := cso.CatSwarmOptimizer{
		CatNum:         *catNum,
		MixtureRatio:   *mr,
		Smp:           	*smp,
		Srd:            *srd,
		VelocityLimit:  *velLimitflag,
		FitnessFunc:    fitnessFunc,
		StateGenerator: stateGenerator,
	}
	

	optimizer.Optimize(*iterations)
}
