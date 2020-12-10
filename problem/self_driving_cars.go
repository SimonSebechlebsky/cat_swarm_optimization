package problem

import (
	"fmt"
	"github.com/SimonSebechlebsky/cat_swarm_optimization/cso"
	"log"
	"math/rand"
	"os"
)

const PermutationCarDelimiter = -1

type Intersection struct {
	Row    int
	Column int
}

type Car struct {
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
	cars := make([]Car, 0, problemDef.VehicleCount)
	cars = append(cars, NewCar())

	for i := 0; i < len(s); i++ {
		if s[i] <= PermutationCarDelimiter {
			curCar++
			cars = append(cars, NewCar())
			continue
		}
		rideIndex := s[i]
		cars[curCar].Rides = append(cars[curCar].Rides, &problemDef.Rides[rideIndex])
	}
	return cars
}


func GetSolutionStateFunc(problemDef *CarProblemDefinition) func() cso.SolutionState {
	
	return func() cso.SolutionState {
		delimiterIdentifier := PermutationCarDelimiter

		perm := rand.Perm(problemDef.RideCount + problemDef.VehicleCount - 1)
		for i, val := range perm {
			if val >= problemDef.RideCount {
				perm[i] = delimiterIdentifier
				delimiterIdentifier -= 1
			}
		}
		s := cso.SolutionState{}
		s = perm
		return s
	}

}

func GetFitnessFunc(problemDef *CarProblemDefinition) func(state cso.SolutionState) int {

	return func(s cso.SolutionState) int {
		fitness := 0
		cars := AssignCars(s, problemDef)
		for _, car := range cars {
			fitness += CarFitness(car, problemDef.Bonus)
		}

		return fitness
	}
}


func CarFitness(car Car, bonus int) int {
	t := 0
	carFitness := 0
	curPos := Intersection{}
	for _, ride := range car.Rides {
		distanceToStart := Distance(curPos, ride.Start)
		rideDistance := Distance(ride.Start, ride.Dest)
		destTime := t + distanceToStart + rideDistance
		if destTime > ride.LatestFinish {
			continue
		}
		carFitness += rideDistance
		if ride.EarliestStart >= t+distanceToStart {
			carFitness += bonus
			waitTime := ride.EarliestStart - (t + distanceToStart)
			t += waitTime
		}
		t += distanceToStart + rideDistance
		curPos = ride.Dest
	}

	return carFitness
}

func Distance(i1 Intersection, i2 Intersection) int {
	return Abs(i1.Column-i2.Column) + Abs(i1.Row-i2.Row)
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
