package cso

import (
	"math"
)

//Array of swaps representing difference between two states
type Velocity []Swap

type Swap struct {
	From int
	To   int
}

//Creates New velocity
func NewVelocity(length int) Velocity {
	return Velocity{}
}

func UpdateState(state SolutionState, velocity Velocity) SolutionState {
	//Adds velocity to SolutionState
	return SolutionState{}
}

func createIndexMap(nums []int) map[int]int {
	indexMap := make(map[int]int, len(nums))
	for i, val := range nums {
		indexMap[val] = i
	}
	return indexMap
}

func smallerFirst(a int, b int) (int, int) {
	if a < b {
		return a, b
	}
	return b, a
}

func getPermutationCycles(perm1 []int, perm2 []int) [][]int {

	if len(perm1) != len(perm2) {
		panic("Permutations need to have same length")
	}
	stateIndexMap := createIndexMap(perm1)

	//create set of numbers in permutation
	numSet := make(map[int]struct{})
	for _, val := range perm1 {
		numSet[val] = struct{}{}
	}

	permutationCycles := make([][]int, 0, 10)

	for key := range numSet {
		permutationCycle := make([]int, 0, 10)
		num := key
		permutationCycle = append(permutationCycle, num)
		for {
			delete(numSet, num)
			stateIndex := stateIndexMap[num]
			mappedNum := perm2[stateIndex]
			num = mappedNum
			if key == mappedNum {
				break
			}
			permutationCycle = append(permutationCycle, num)
		}

		if len(permutationCycle) > 1 {
			permutationCycles = append(permutationCycles, permutationCycle)
		}
	}
	return permutationCycles
}

// Computes velocity in relationship to another state -
// swaps that need to be made to get to that state
func (state SolutionState) GetVelocity(state2 SolutionState) Velocity {
	permCycles := getPermutationCycles(state, state2)
	swaps := make([]Swap, 0, 10)

	for _, cycle := range permCycles {
		for _, num := range cycle[1:] {
			from, to := smallerFirst(cycle[0], num)
			swaps = append(swaps, Swap{from, to})
		}
	}
	return swaps
}

//Applies Array of swaps to SolutionState
func (state SolutionState) ApplyVelocity(velocity Velocity) SolutionState {
	finalState := SolutionState(make([]int, len(state)))
	copy(finalState, state)
	indexMap := createIndexMap(finalState)

	for _, swap := range velocity {
		indexFrom := indexMap[swap.From]
		indexTo := indexMap[swap.To]
		indexMap[swap.From] = indexTo
		indexMap[swap.To] = indexFrom
		finalState[indexFrom], finalState[indexTo] =
			finalState[indexTo], finalState[indexFrom]
	}
	return finalState
}

//Float f needs to be positive
func (velocity Velocity) MultiplyByFloat(f float64) Velocity {
	intPart := int(math.Floor(f))
	floatPart := f - float64(intPart)

	newVelocity := velocity.Repeat(intPart)
	newVelocity = append(newVelocity, velocity.Shrink(floatPart)...)

	return newVelocity
}

//Float f needs to be < 1
func (velocity Velocity) Shrink(f float64) Velocity {
	swapCount := int(math.Ceil(float64(len(velocity)) * f))
	velocity = velocity[:swapCount]
	return velocity
}

// Outputs velocity which is equivalent to original velocity repeated n times
func (velocity Velocity) Repeat(n int) Velocity {
	if n == 0 {
		return Velocity([]Swap(nil))
	}

	newVelocity := Velocity(make([]Swap, 0, len(velocity)*(n+1)))
	for i := 0; i < n; i++ {
		newVelocity = append(newVelocity, velocity...)
	}
	return newVelocity
}

//Outputs equivalent velocity with minimal length
func (velocity Velocity) Minimize(stateGenerator func() SolutionState) Velocity {
	state := stateGenerator()
	finalState := state.ApplyVelocity(velocity)
	minimizedVelocity := state.GetVelocity(finalState)
	return minimizedVelocity
}

func (velocity Velocity) Add(velocity2 Velocity, stateGenerator func() SolutionState) Velocity {
	mergedSwaps := append([]Swap(nil), velocity...)
	mergedSwaps = append(mergedSwaps, velocity2...)
	finalVelocity := Velocity(mergedSwaps)
	return finalVelocity.Minimize(stateGenerator)
}
