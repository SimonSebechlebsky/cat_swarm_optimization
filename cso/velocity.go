package cso

//Array of swaps representing difference between two states
type Velocity struct {
	Swaps []Swap
}

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
	indexMap := make(map[int]int)
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
func (state *SolutionState) GetVelocity(state2 SolutionState) Velocity {
	permCycles := getPermutationCycles(state.Permutation, state2.Permutation)
	swaps := make([]Swap, 0, 10)

	for _, cycle := range permCycles {
		for _, num := range cycle[1:] {
			from, to := smallerFirst(cycle[0], num)
			swaps = append(swaps, Swap{from, to})
		}
	}
	return Velocity{Swaps: swaps}
}

//Apllies Array of swaps to SolutionState
func (state *SolutionState) ApplyVelocity(velocity Velocity) SolutionState {
	finalState := SolutionState{make([]int, len(state.Permutation))}
	copy(finalState.Permutation, state.Permutation)
	indexMap := createIndexMap(finalState.Permutation)

	for _, swap := range velocity.Swaps {
		indexFrom := indexMap[swap.From]
		indexTo := indexMap[swap.To]
		indexMap[swap.From] = indexTo
		indexMap[swap.To] = indexFrom
		finalState.Permutation[indexFrom], finalState.Permutation[indexTo] =
			finalState.Permutation[indexTo], finalState.Permutation[indexFrom]
	}
	return finalState
}

func (velocity Velocity) MultiplyByFloat(f float32) Velocity {
	// Multiplies velocity by float, defined in Discrete Particle Swarm Optimization, illustrated by the Traveling Salesman Problem
	return Velocity{}
}

func (velocity Velocity) Add(velocity2 Velocity) Velocity {
	// Get minimal velocity which is equal to velocity+velocity2
	return Velocity{}
}
