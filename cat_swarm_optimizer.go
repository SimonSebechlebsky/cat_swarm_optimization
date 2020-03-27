package main

type Mode int

const (
	tracingMode Mode = iota
	seekingMode
)

type CatSwarmOptimizer struct {
	catNum       int
	mixtureRatio float32 // ratio of cats in tracing and seeking mode
	smp          int     // seeking memory pool - how many cat copies to spawn in seeking mode
	srd          float32 // seeking range of the selected dimension - maximum range to mutate a dimension
	cdc          int     //counts of dimension to change - how many dimensions will be varied
	cats         []Cat
}

type Cat struct {
	mode Mode
	SolutionState
}

type SolutionState struct {
}

func Fitness(s SolutionState) int {
	return 0
}
