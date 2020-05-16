package cso

import (
	wr "github.com/mroth/weightedrand"
	"math"
	"math/rand"
	"time"
)

type Mode int

const (
	TracingMode Mode = iota
	SeekingMode
)

type ChanResult struct {
    Id int
    Value int
}

type CatSwarmOptimizer struct {
	CatNum        int
	MixtureRatio  float32 // ratio of cats in tracing and seeking mode
	Smp           int     // seeking memory pool - how many cat copies to spawn in seeking mode
	Srd           int     // seeking range of the selected dimension - maximum number of permutations in seeking mode
	VelocityLimit int     // tracing mode velocity limit
	Cats          []Cat
	FitnessFunc   func(state SolutionState) float32
	BestCat       Cat
}

type Cat struct {
	Mode     Mode
	State    SolutionState
	Vel Velocity
	FitnessFunc   func(state SolutionState) int
}

type SolutionState struct {
	Permutation []int
}

func GetRandomVelocity(state SolutionState, srd int) Velocity {
	rand.Seed(time.Now().UnixNano())
	randomSwapsCount := rand.Intn(srd)
	swaps := make([]Swap, 0, randomSwapsCount)

	for i := 0; i < randomSwapsCount; i++ {
		from, to := rand.Intn(len(state.Permutation)), rand.Intn(len(state.Permutation))
		for from == to {
			from, to = rand.Intn(len(state.Permutation)), rand.Intn(len(state.Permutation))
		}
		swaps = append(swaps, Swap{state.Permutation[from], state.Permutation[to]})
	}
	return Velocity{Swaps: swaps}
}

func (state SolutionState) MoveStateRandomly(srd int) SolutionState {
	// Move state randomly of max srd permutations
	RandomState := SolutionState{make([]int, len(state.Permutation))}
	velocity := GetRandomVelocity(state, srd)

    RandomState = state.ApplyVelocity(velocity)

	return RandomState
}

func GetProbability(minFitness, maxFitness, iFitness int) float64 {
	prob := math.Abs(float64(iFitness - minFitness)) / float64(maxFitness - minFitness)
    return prob
}

func FindMin(arr []int) (min int) {
    for i, value := range arr {
        if i == 0 || value < min {
            min = value
        }
    }
    return
}

func FindMax(arr []int) (max int) {
    for i, value := range arr {
        if i == 0 || value > max {
            max = value
        }
    }
    return
}

func FitnessFuncWrapper(state SolutionState, id int, fitnessChan chan ChanResult, fitnessFunc func(state SolutionState) int) {
    fitness := fitnessFunc(state)
    result := ChanResult{Id: id, Value: fitness}
	fitnessChan <- result
}

func (cat *Cat) Seek(smp int, srd int) {
	// Do all the steps of seeking mode and set new SolutionState to cat
	rand.Seed(time.Now().UnixNano())
    stateCopies := make([]SolutionState, 0, smp)
    fitnessArray := make([]int, smp)
    for i := 0; i < smp; i++ {
        state := cat.State.MoveStateRandomly(srd)
        stateCopies = append(stateCopies, state)
    }
    fitnessChan := make(chan ChanResult, len(stateCopies))
    for i, state := range stateCopies {
        go FitnessFuncWrapper(state, i, fitnessChan, cat.FitnessFunc)
    }

    for i := 0; i < len(stateCopies); i++{
        fitness := <-fitnessChan
        fitnessArray[fitness.Id] = fitness.Value
    }

    minFitness := FindMin(fitnessArray)
    maxFitness := FindMax(fitnessArray)
    weightedCats := make([]wr.Choice, 0, len(stateCopies))
    for i, state := range stateCopies {
        probability := GetProbability(minFitness, maxFitness, fitnessArray[i])
        probInt := uint(probability*1000)
        weightedCats = append(weightedCats, wr.Choice{Item: state, Weight: probInt})
    }

    c := wr.NewChooser(weightedCats...)
    result := c.Pick().(SolutionState)

    cat.State = result
}

func (cat *Cat) Trace() {
	// Do all the steps of tracing mode and set new SolutionState to cat
}

func (cat Cat) UpdateVelocity() Velocity {
	// Velocity update formula from tracing mode (step1 in article)
	return Velocity{}
}

func (optimizer CatSwarmOptimizer) Optimize(steps int) SolutionState {

	return SolutionState{}
}
