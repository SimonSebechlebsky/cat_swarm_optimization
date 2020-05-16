package cso_test

import (
    "fmt"
    "github.com/SimonSebechlebsky/cat_swarm_optimization/cso"
    "github.com/SimonSebechlebsky/cat_swarm_optimization/problem"
    "testing"
)

func TestSeek(t *testing.T) {
    problemDef := problem.LoadCarProblemDefinition("../problem/inputs/b_should_be_easy.in")
    FitnessFunc := problem.GetFitnessFunc(problemDef)
    SolutionStateFunc := problem.GetSolutionStateFunc(problemDef)
    state := SolutionStateFunc()
    fmt.Println("state: ", state)

    cat := cso.Cat{Mode: cso.SeekingMode, State: state, Vel: cso.Velocity{}, FitnessFunc: FitnessFunc}
    fitness := FitnessFunc(state)

    for i := 0; i < 1000; i++ {
        cat.Seek(200, 5)
        fitness = FitnessFunc(cat.State)
        fmt.Println(fitness, cat.State)
    }

}