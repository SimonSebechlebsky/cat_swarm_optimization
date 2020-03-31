package main

import (
	"fmt"
	"log"
	"os"
)

type Intersection struct {
	Row    int
	Column int
}

type Ride struct {
	Start         Intersection
	Dest          Intersection
	EarliestStart int
	LatestFinish  int
}

type ProblemDefinition struct {
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

func (problemDef ProblemDefinition) String() string {
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

func LoadProblemDefinition(InputFilePath string) *ProblemDefinition {
	log.Print(fmt.Sprintf("Opening file %s", InputFilePath))
	inputFile, err := os.Open(InputFilePath)
	if err != nil {
		panic(err)
	}

	problemDef := ProblemDefinition{}
	problemDef.FileName = InputFilePath
	_, _ = fmt.Fscanf(
		inputFile, "%d %d %d %d %d %d\n",
		&problemDef.Rows, &problemDef.Columns, &problemDef.VehicleCount,
		&problemDef.RideCount, &problemDef.Bonus, &problemDef.TotalTimeSteps)

	for i := 0; i < problemDef.RideCount; i++ {
		ride := NewRide()
		fmt.Fscanf(inputFile, "%d %d %d %d %d %d\n",
			&ride.Start.Row, &ride.Start.Column, &ride.Dest.Row, &ride.Dest.Column,
			&ride.EarliestStart, &ride.LatestFinish)
		problemDef.Rides = append(problemDef.Rides, ride)
	}
	log.Print(problemDef)
	return &problemDef
}
