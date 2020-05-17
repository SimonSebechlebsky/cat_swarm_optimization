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
	Id    int
	Value int
}

type CatSwarmOptimizer struct {
	CatNum         int
	MixtureRatio   float64 // ratio of cats in tracing and seeking mode
	Smp            int     // seeking memory pool - how many cat copies to spawn in seeking mode
	Srd            int     // seeking range of the selected dimension - maximum number of permutations in seeking mode
	VelocityLimit  int     // tracing mode velocity limit
	Cats           []Cat
	FitnessFunc    func(state SolutionState) int
	StateGenerator func() SolutionState
	BestCat        *Cat
}

type Cat struct {
	Mode           Mode
	State          SolutionState
	Vel            Velocity
	VelocityLimit  int
	FitnessFunc    func(state SolutionState) int
	StateGenerator func() SolutionState
	Fitness 	   int
}

type SolutionState struct {
	Permutation []int
}

func GetRandomVelocity(state SolutionState, srd int) Velocity {
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
	prob := math.Abs(float64(iFitness-minFitness)) / float64(maxFitness-minFitness)
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

func FitnessFuncWrapper(state SolutionState, id int, fitnessChan chan ChanResult,
	fitnessFunc func(state SolutionState) int) {
	fitness := fitnessFunc(state)
	result := ChanResult{Id: id, Value: fitness}
	fitnessChan <- result
}

func (cat *Cat) Seek(smp int, srd int) {
	// Do all the steps of seeking mode and set new SolutionState to cat
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

	for i := 0; i < len(stateCopies); i++ {
		fitness := <-fitnessChan
		fitnessArray[fitness.Id] = fitness.Value
	}

	minFitness := FindMin(fitnessArray)
	maxFitness := FindMax(fitnessArray)
	weightedCats := make([]wr.Choice, 0, len(stateCopies))
	for i, _ := range stateCopies {
		probability := GetProbability(minFitness, maxFitness, fitnessArray[i])
		probInt := uint(probability * 1000)
		weightedCats = append(weightedCats, wr.Choice{Item: i, Weight: probInt})
	}

	c := wr.NewChooser(weightedCats...)
	resultId := c.Pick().(int)

	cat.State = stateCopies[resultId]
	cat.Fitness = fitnessArray[resultId]
}

func (cat *Cat) Trace(bestCat *Cat) {
	cat.UpdateVelocity(bestCat.State)
	cat.State = cat.State.ApplyVelocity(cat.Vel)
	cat.Fitness = cat.FitnessFunc(cat.State)
}

func (cat *Cat) UpdateVelocity(bestState SolutionState) {
	r := rand.Float64()
	c := 1. // TODO maybe move c to optimizer configuration params? There's nothing in the paper about it though
	bestVelocity := cat.State.GetVelocity(bestState)
	modifiedVelocity := bestVelocity.MultiplyByFloat(r * c)
	newVelocity := cat.Vel.Add(modifiedVelocity, cat.StateGenerator)

	if len(newVelocity.Swaps) > cat.VelocityLimit {
		newVelocity.Swaps = newVelocity.Swaps[:cat.VelocityLimit]
	}
	cat.Vel = newVelocity
}

func (optimizer CatSwarmOptimizer) CreateCats() []Cat {
	cats := make([]Cat, 0, optimizer.CatNum)
	for i := 0; i < optimizer.CatNum; i++ {
		state := optimizer.StateGenerator()
		velocity := GetRandomVelocity(state, optimizer.VelocityLimit)
		fitness := optimizer.FitnessFunc(state)

		cats = append(cats, Cat{Mode: nil, State: state, Vel: velocity, VelocityLimit: optimizer.VelocityLimit,
			FitnessFunc: optimizer.FitnessFunc, StateGenerator: optimizer.StateGenerator, Fitness: fitness})
	}

	return cats
}

func SetMode(cats []Cat, mixtureRatio float64) []Cat {
	tracingCatsCount := int(math.Round(mixtureRatio * float64(len(cats))))
	for i, cat := range cats {
		if i < tracingCatsCount {
			cat.Mode = TracingMode
		} else {
			cat.Mode = SeekingMode
		}
	}
	rand.Shuffle(len(cats), func(i, j int) {
		cats[i], cats[j] = cats[j], cats[i]
	})
	return cats
}

func FindBestCat(cats []Cat) *Cat {
	bestCat := cats[0]
	for _, cat := range cats {
		if cat.Fitness > bestCat.Fitness {
			bestCat = cat
		}
	}
	return &bestCat
}

func (optimizer CatSwarmOptimizer) Optimize(num int) SolutionState {
	rand.Seed(time.Now().UnixNano())
	cats := optimizer.CreateCats()

	for i := 0; i < num; i++ {
		cats = SetMode(cats, optimizer.MixtureRatio)

		optimizer.BestCat = FindBestCat(cats)

		for _, cat := range cats {
			if cat.Mode == SeekingMode {
				cat.Seek(optimizer.Smp, optimizer.Srd)
			} else {
				cat.Trace(optimizer.BestCat)
			}
		}
	}

	optimizer.BestCat = FindBestCat(cats)
	return optimizer.BestCat.State
}
