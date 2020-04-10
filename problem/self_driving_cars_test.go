package problem

import (
	"github.com/SimonSebechlebsky/cat_swarm_optimization/cso"
	"testing"
)

func TestAssignCars(t *testing.T) {
	problemDef := LoadCarProblemDefinition("./inputs/a_example.in")
	solState := cso.SolutionState{Permutation: []int{1, 2, PermutationCarDelimiter, 0}}
	cars := AssignCars(solState, problemDef)
	if len(cars) != 2 {
		t.Errorf("Expected 2 cars, got %d", len(cars))
	}

	if len(cars[0].Rides) != 2 {
		t.Errorf("Expected 2 rides for car 0, got %d", len(cars[0].Rides))
	}

	if len(cars[1].Rides) != 1 {
		t.Errorf("Expected 1 ride for car 1, got %d", len(cars[1].Rides))
	}

}

func TestCarFitness1(t *testing.T) {
	problemDef := LoadCarProblemDefinition("./inputs/a_example.in")
	solState := cso.SolutionState{Permutation: []int{0, PermutationCarDelimiter, 2, 1}}
	cars := AssignCars(solState, problemDef)
	fitness := CarFitness(cars[0], problemDef.Bonus)

	if fitness != 6 {
		t.Errorf("Expected 6 fitness points, got %d", fitness)
	}
}

func TestCarFitness2(t *testing.T) {
	problemDef := LoadCarProblemDefinition("./inputs/a_example.in")
	solState := cso.SolutionState{Permutation: []int{0, PermutationCarDelimiter, 2, 1}}
	cars := AssignCars(solState, problemDef)
	fitness := CarFitness(cars[1], problemDef.Bonus)

	if fitness != 4 {
		t.Errorf("Expected 6 fitness points, got %d", fitness)
	}
}

func TestCarFitnessSkipRide(t *testing.T) {
	car := Car{
		Rides: []*Ride{
			{
				Start: Intersection{
					Row:    2,
					Column: 2,
				},
				Dest: Intersection{
					Row:    5,
					Column: 7,
				},
				EarliestStart: 0,
				LatestFinish:  10, //should skip this ride
			},
			{
				Start: Intersection{
					Row:    0,
					Column: 0,
				},
				Dest: Intersection{
					Row:    5,
					Column: 5,
				},
				EarliestStart: 0,
				LatestFinish:  10,
			},
		},
	}
	fitness := CarFitness(car, 0)

	if fitness != 10 {
		t.Errorf("Expected 10 fitness points, got %d", fitness)
	}
}

func TestFitness(t *testing.T) {
	problemDef := LoadCarProblemDefinition("./inputs/a_example.in")
	solState := cso.SolutionState{Permutation: []int{0, PermutationCarDelimiter, 2, 1}}
	FitnessFunc := GetFitnessFunc(problemDef)
	fitness := FitnessFunc(solState)
	if fitness != 10 {
		t.Errorf("Expected 10 fitness points, got %d", fitness)
	}
}
