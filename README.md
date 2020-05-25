## About
This repository is implementation of Cat swarm optimization algorithm modified for discrete permutational problems. It contains 
also implementation of the hashcode 2018 problem (self driving rides) with the fitness function for this optimizer.

## Requirements
go>=1.14

## Build & Run
```
go build
./cat_swarm_optimization 
```
You can specify  input files and hyperparameters via commandline arguments - e.g. `./cat_swarm_optimization --input_file ./problem/inputs/c_no_hurry_in --cat_num 50`
View the `main.go` file for the list of possible parameters. 
