package cso

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetPermutationCycles(t *testing.T) {
	permCycles := getPermutationCycles([]int{1, 2, 3, 4, 5}, []int{2, 3, 1, 4, 5})

	if len(permCycles) != 1 {
		t.Errorf("Expected 2 permutation cycles, got %d", len(permCycles))
	}

	if !assert.Equal(t, permCycles[0], []int{1, 2, 3}) {
		t.Errorf("Expected permutation cycle [1,2,3], got %v", permCycles[0])
	}
}

func TestGetVelocity(t *testing.T) {
	state1 := SolutionState{[]int{1, 2, 3, 4, 5, 6, 7}}
	state2 := SolutionState{[]int{2, 1, 5, 4, 6, 3, 7}}
	velocity := state1.GetVelocity(state2) // Permutation cycles (1 2) (3 5 6), 3 swaps (i.e (1 2) (3 5) (3 6))

	finalState := state1.ApplyVelocity(velocity)
	assert.Equal(t, state2, finalState)
}

func TestMultiplyByFloatLessThanOne(t *testing.T) {
	velocity := Velocity{[]Swap{{1, 2}, {2, 5}, {6, 8}, {4, 7}, {6, 3}}}
	newVelocity := velocity.MultiplyByFloat(0.4)

	assert.Equal(t, newVelocity.Swaps, []Swap{{1, 2}, {2, 5}})
}

func TestMultiplyByFloatMoreThanOne(t *testing.T) {
	velocity := Velocity{[]Swap{{1, 2}, {2, 5}, {6, 8}}}
	newVelocity := velocity.MultiplyByFloat(2.2)

	assert.Equal(t, newVelocity.Swaps, []Swap{{1, 2}, {2, 5}, {6, 8}, {1, 2}, {2, 5}, {6, 8}, {1, 2}})
}

func mockSolutionGenerator() SolutionState {
	return SolutionState{Permutation: []int{1,2,3,4,5,6,7,8}}
}

func TestMinimizeVelocity(t *testing.T) {
	velocity := Velocity{[]Swap{{1, 2}, {2,1}, {1,2}, {2, 5}, {6, 8}}}
	minimizedVelocity := velocity.Minimize(mockSolutionGenerator)
	expectedState := SolutionState{Permutation: []int{5,1,3,4,2,8,7,6}}
	startState := mockSolutionGenerator()
	destState := startState.ApplyVelocity(minimizedVelocity)

	assert.Equal(t, len(minimizedVelocity.Swaps), 3)
	assert.Equal(t, destState.Permutation, expectedState.Permutation)
}


func TestAddVelocity(t *testing.T) {
	velocity1 := Velocity{[]Swap{{1, 2}, {2, 5}, {6, 8}}}
	velocity2 := Velocity{[]Swap{{3, 4}, {3, 7}, {7, 3}}}
	minimizedVelocity := velocity1.Add(velocity2, mockSolutionGenerator)
	startState := mockSolutionGenerator()
	destState := startState.ApplyVelocity(minimizedVelocity)
	expectedState := startState.ApplyVelocity(velocity1).ApplyVelocity(velocity2)

	assert.Equal(t, len(minimizedVelocity.Swaps), 4)
	assert.Equal(t, destState.Permutation, expectedState.Permutation)
}