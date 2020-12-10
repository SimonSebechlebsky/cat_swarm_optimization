package cso_test

import (
    "fmt"
    "github.com/SimonSebechlebsky/cat_swarm_optimization/cso"
    "github.com/SimonSebechlebsky/cat_swarm_optimization/problem"
    "io/ioutil"
    "log"
    "os"
    "testing"
)


func TestMain(m *testing.M) {
    log.SetOutput(ioutil.Discard)
    os.Exit(m.Run())
}


func TestSeek(t *testing.T) {
    problemDef := problem.LoadCarProblemDefinition("../problem/inputs/b_should_be_easy.in")
    FitnessFunc := problem.GetFitnessFunc(problemDef)
    SolutionStateFunc := problem.GetSolutionStateFunc(problemDef)
    state := SolutionStateFunc()
    fmt.Println("state: ", state)

    cat := cso.Cat{Mode: cso.SeekingMode, State: state, Vel: cso.Velocity{}, FitnessFunc: FitnessFunc}
    fitness := FitnessFunc(state)

    for i := 0; i < 200; i++ {
        cat.Seek(50, 3)
        fitness = FitnessFunc(cat.State)
        fmt.Println(fitness, cat.State)
    }

}


func BenchmarkOptimization(b *testing.B) {

    problemDef := problem.LoadCarProblemDefinition("../problem/inputs/b_should_be_easy.in")
    fitnessFunc := problem.GetFitnessFunc(problemDef)
    stateGenerator := problem.GetSolutionStateFunc(problemDef)

    optimizer := cso.CatSwarmOptimizer{
        CatNum:         20,
        MixtureRatio:   0.7,
        Smp:           	50,
        Srd:            5,
        VelocityLimit:  10,
        FitnessFunc:    fitnessFunc,
        StateGenerator: stateGenerator,
    }

    for i := 0; i < b.N; i++ {
        optimizer.Optimize(1000)
    }

}


//func TestTrace(t *testing.T) {
//    rand.Seed(time.Now().UnixNano())
//    problemDef := problem.LoadCarProblemDefinition("../problem/inputs/b_should_be_easy.in")
//    FitnessFunc := problem.GetFitnessFunc(problemDef)
//    SolutionStateFunc := problem.GetSolutionStateFunc(problemDef)
//
//    cat := cso.Cat{
//        Mode: cso.TracingMode,
//        State: SolutionStateFunc(),
//        Vel: cso.Velocity{Swaps: []cso.Swap(nil)},
//        VelocityLimit: 5,
//        FitnessFunc: FitnessFunc,
//        StateGenerator:SolutionStateFunc,
//    }
//
//    otherCat := cso.Cat{
//        Mode: cso.TracingMode,
//        State: SolutionStateFunc(),
//        Vel: cso.Velocity{Swaps: []cso.Swap(nil)},
//        FitnessFunc: FitnessFunc,
//        StateGenerator: SolutionStateFunc,
//    }
//
//    for i := 0; i < 1000; i++ {
//        cat.Trace(&otherCat)
//        fitnessCurrent := FitnessFunc(cat.State)
//        fitnessOther := FitnessFunc(otherCat.State)
//        fmt.Println(fitnessCurrent, fitnessOther) // the cat should move towards the other one
//    }
//
//}