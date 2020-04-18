package cso

type Mode int

const (
	tracingMode Mode = iota
	seekingMode
)

type CatSwarmOptimizer struct {
	catNum        int
	mixtureRatio  float32 // ratio of cats in tracing and seeking mode
	smp           int     // seeking memory pool - how many cat copies to spawn in seeking mode
	srd           int     // seeking range of the selected dimension - maximum number of permutations in seeking mode
	velocityLimit int     // tracing mode velocity limit
	cats          []Cat
}

type Cat struct {
	mode     Mode
	state    SolutionState
	velocity Velocity
}

type SolutionState struct {
	Permutation []int
}

func (state SolutionState) MoveStateRandomly(srd int) SolutionState {
	// Move state randomly of max srd permutations
	return SolutionState{}
}

func (cat *Cat) Seek() {
	// Do all the steps of seeking mode and set new SolutionState to cat
}

func (cat *Cat) Trace() {
	// Do all the steps of tracing mode and set new SolutionState to cat
}

func (cat Cat) UpdateVelocity() Velocity {
	// Velocity update formula from tracing mode (step1 in arcitcle)
	return Velocity{}
}

func (optimizer CatSwarmOptimizer) Optimize(steps int) SolutionState {

	return SolutionState{}
}
