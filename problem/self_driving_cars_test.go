package problem

import (
	"github.com/SimonSebechlebsky/cat_swarm_optimization/cso"
	"math/rand"
	"testing"
)


func TestAssignCars(t *testing.T) {
	rand.Seed(0)
	problemDef := LoadCarProblemDefinition("./inputs/a_example.in")
	solState := cso.SolutionState{Permutation: []int {1, 2, PermutationCarDelimiter, 0}}
	cars := AssignCars(solState, problemDef)
	if len(cars) != 2 {
		t.Errorf("Expected 2 cars, got %d",  len(cars))
	}

	if len(cars[0].Rides) != 2 {
		t.Errorf("Expected 2 rides for car 0, got %d",  len(cars[0].Rides))
	}

	if len(cars[1].Rides) != 1 {
		t.Errorf("Expected 1 ride for car 1, got %d",  len(cars[1].Rides))
	}

}