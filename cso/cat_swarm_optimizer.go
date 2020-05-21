package cso

import (
	wr "github.com/mroth/weightedrand"
	"log"
	"math"
	"math/rand"
	"sync"
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
	FitnessFunc    func(state SolutionState) int
	StateGenerator func() SolutionState
	BestCat        Cat
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

type SolutionState []int

func GetRandomVelocity(state SolutionState, srd int) Velocity {
	randomSwapsCount := rand.Intn(srd)
	swaps := make([]Swap, 0, randomSwapsCount)

	for i := 0; i < randomSwapsCount; i++ {
		from, to := rand.Intn(len(state)), rand.Intn(len(state))
		for from == to {
			from, to = rand.Intn(len(state)), rand.Intn(len(state))
		}
		swaps = append(swaps, Swap{state[from], state[to]})
	}
	return swaps
}

func (state SolutionState) MoveStateRandomly(srd int) SolutionState {
	// Move state randomly of max srd permutations
	RandomState := SolutionState(make([]int, len(state)))
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

	if len(newVelocity) > cat.VelocityLimit {
		newVelocity = newVelocity[:cat.VelocityLimit]
	}
	cat.Vel = newVelocity
}

func (optimizer CatSwarmOptimizer) CreateCats() []Cat {
	cats := make([]Cat, 0, optimizer.CatNum)
	for i := 0; i < optimizer.CatNum; i++ {
		state := optimizer.StateGenerator()
		velocity := GetRandomVelocity(state, optimizer.VelocityLimit)
		fitness := optimizer.FitnessFunc(state)

		cats = append(cats, Cat{Mode: SeekingMode, State: state, Vel: velocity, VelocityLimit: optimizer.VelocityLimit,
			FitnessFunc: optimizer.FitnessFunc, StateGenerator: optimizer.StateGenerator, Fitness: fitness})
	}

	return cats
}

func SetMode(cats []Cat, mixtureRatio float64) []Cat {
	tracingCatsCount := int(math.Round(mixtureRatio * float64(len(cats))))
	for i, _ := range cats {
		if i < tracingCatsCount {
			cats[i].Mode = TracingMode
		} else {
			cats[i].Mode = SeekingMode
		}
	}
	rand.Shuffle(len(cats), func(i, j int) {
		cats[i], cats[j] = cats[j], cats[i]
	})
	return cats
}

func FindBestCat(cats []Cat, bestCat *Cat) Cat {
	for i, _ := range cats {
		if cats[i].Fitness > bestCat.Fitness {
			bestCat = &cats[i]
		}
	}
	bestCatCopy := *bestCat
	bestCatCopy.State = SolutionState(nil)
	bestCatCopy.State = append(bestCatCopy.State, bestCat.State...)
	return bestCatCopy
}

func (optimizer CatSwarmOptimizer) Optimize(num int) SolutionState {
	rand.Seed(time.Now().UnixNano())
	wg := sync.WaitGroup{}
	cats := optimizer.CreateCats()
	optimizer.BestCat = Cat{}

	for i := 0; i < num; i++ {
		cats = SetMode(cats, optimizer.MixtureRatio)
		fitnessSum := 0
		optimizer.BestCat = FindBestCat(cats, &optimizer.BestCat)

		for i, _ := range cats {
			wg.Add(1)
			cat := &cats[i]

			go func(cat *Cat, optimizer *CatSwarmOptimizer, wg *sync.WaitGroup) {
				defer wg.Done()

				if cat.Mode == SeekingMode {
					cat.Seek(optimizer.Smp, optimizer.Srd)
				} else {
					cat.Trace(&optimizer.BestCat)
				}
			}(cat, &optimizer, &wg)
		}
		wg.Wait()

		for i, _ := range cats {
			fitnessSum += cats[i].Fitness
		}

		log.Printf("Highest fitness in iteration %d: %d, average fitness in iteration %d: %f",
			i+1, optimizer.BestCat.Fitness, i+1, float32(fitnessSum)/float32(len(cats)))
	}

	optimizer.BestCat = FindBestCat(cats, &optimizer.BestCat)
	return optimizer.BestCat.State
}
