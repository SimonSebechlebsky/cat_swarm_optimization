package problem

import (
	"fmt"
	"github.com/SimonSebechlebsky/cat_swarm_optimization/cso"
	"log"
	"math/rand"
	"os"
)

const PermutationCarDelimiter  = -1

type Intersection struct {
	Row    int
	Column int
}

type Car struct {
	CurPos Intersection
	Rides []*Ride
}

type Ride struct {
	Start         Intersection
	Dest          Intersection
	EarliestStart int
	LatestFinish  int
}

type CarProblemDefinition struct {
	FileName       string
	Rows           int
	Columns        int
	VehicleCount   int
	RideCount      int
	Bonus          int //Added if ride starts at earliest possible moment
	TotalTimeSteps int
	Rides          []Ride
}

func NewRide() Ride {
	ride := Ride{}
	ride.Start = Intersection{}
	ride.Dest = Intersection{}
	return ride
}

func NewCar() Car {
	car := Car{}
	car.CurPos = Intersection{}
	return car
}

func (problemDef CarProblemDefinition) String() string {
	stringFormat := `
InputFile: %s,
Rows: %d, Columns: %d,
VehicleCount: %d, RideCount %d,
Bonus: %d,
TotalTimeSteps: %d
`
	return fmt.Sprintf(stringFormat, problemDef.FileName, problemDef.Rows, problemDef.Columns,
		problemDef.VehicleCount, problemDef.RideCount, problemDef.Bonus, problemDef.TotalTimeSteps)
}

func LoadCarProblemDefinition(InputFilePath string) *CarProblemDefinition {
	log.Print(fmt.Sprintf("Opening file %s", InputFilePath))
	inputFile, err := os.Open(InputFilePath)
	if err != nil {
		panic(err)
	}

	problemDef := CarProblemDefinition{}
	problemDef.FileName = InputFilePath
	_, _ = fmt.Fscanf(
		inputFile, "%d %d %d %d %d %d\n",
		&problemDef.Rows, &problemDef.Columns, &problemDef.VehicleCount,
		&problemDef.RideCount, &problemDef.Bonus, &problemDef.TotalTimeSteps)

	for i := 0; i < problemDef.RideCount; i++ {
		ride := NewRide()
		_, _ = fmt.Fscanf(inputFile, "%d %d %d %d %d %d\n",
			&ride.Start.Row, &ride.Start.Column, &ride.Dest.Row, &ride.Dest.Column,
			&ride.EarliestStart, &ride.LatestFinish)
		problemDef.Rides = append(problemDef.Rides, ride)
	}
	log.Print(problemDef)
	return &problemDef
}


func AssignCars(s cso.SolutionState, problemDef *CarProblemDefinition) []Car {
	curCar := 0
	var cars []Car
	cars = append(cars, NewCar())

	for i:= 0; i < len(s.Permutation); i++ {
		if	s.Permutation[i] == PermutationCarDelimiter {
			curCar++
			cars = append(cars, NewCar())
			continue
		}
		rideIndex := s.Permutation[i]
		cars[curCar].Rides = append(cars[curCar].Rides, &problemDef.Rides[rideIndex])
	}
	return cars
}


func NewSolutionState(rideCount int, vehicleCount int) cso.SolutionState {

	perm := rand.Perm(rideCount+vehicleCount-1)
	for i, val := range perm {
		if val >= rideCount {
			perm[i] = PermutationCarDelimiter
		}
	}
	s := cso.SolutionState{}
	s.Permutation = perm
	return s
}


func Fitness(s cso.SolutionState, problemDef *CarProblemDefinition) int {
	fitness := 0
	cars := AssignCars(s, problemDef)
	fitnessChan := make(chan int, len(cars))
	for _, car := range cars {
		go func() {
			carFitness := CarFitness(car, problemDef.Bonus)
			fitnessChan <- carFitness
		}()
	}

	for i:=0; i < len(cars); i++ {
		fitness += <-fitnessChan
	}
	return fitness
}


func CarFitness(car Car, bonus int) int {
	t := 0
	carFitness := 0
	for _, ride := range car.Rides {
		distanceToStart := Distance(car.CurPos, ride.Start)
		rideDistance := Distance(ride.Start, ride.Dest)
		destTime := t+distanceToStart+rideDistance
		if destTime > ride.LatestFinish {
			continue
		}
		carFitness += rideDistance
		if ride.EarliestStart >= t+distanceToStart {
			carFitness += bonus
		}
		//TODO rest of fitness
	}

	return carFitness
}


func Distance(i1 Intersection, i2 Intersection) int {
	return Abs(i1.Column-i2.Column)+Abs(i1.Row-i2.Row)
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

