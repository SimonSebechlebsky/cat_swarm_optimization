package cso

type Velocity struct {
	exchangeNodes []ExchangeNode
}

type ExchangeNode struct {
	From int
	To int
}

func NewVelocity(length int) Velocity{
	//Creates New velocity
	return Velocity{}
}

func UpdateState(state SolutionState, velocity Velocity) SolutionState {
	//Adds velocity to SolutionState
	return SolutionState{}
}

func GetVelocity(s1 SolutionState, s2 SolutionState) Velocity {
	//Calulates difference between 2 states, defined in Discrete Particle Swarm Optimization, illustrated by the Traveling Salesman Problem
	return Velocity{}
}


func (velocity Velocity) MultiplyByFloat(f float32) Velocity {
	// Multiplies velocity by float, defined in Discrete Particle Swarm Optimization, illustrated by the Traveling Salesman Problem
	return Velocity{}
}


func (velocity Velocity) Add(velocity2 Velocity) Velocity {
	// Get minimal velocity which is equal to velocity+velocity2
	return Velocity{}
}


