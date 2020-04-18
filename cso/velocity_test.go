package cso

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestGetPermutationCycles(t *testing.T) {
	permCycles := getPermutationCycles([]int{1, 2, 3, 4, 5}, []int{2, 3, 1, 4, 5})

	if len(permCycles) != 1 {
		t.Errorf("Expected 2 permutation cycles, got %d", len(permCycles))
	}

	if !reflect.DeepEqual(permCycles[0], []int{1, 2, 3}) {
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
