package simulation

import (
	"fmt"
	"github.com/kinethic"
)

func SimulateKineticModel(initialConditions map[string]float64, params map[string]float64, temperatureProfile []float64, timeSteps int, dt float64) [][]float64 {
	// Extract initial conditions
	X := initialConditions["X"]
	P := initialConditions["P"]
	S := initialConditions["S"]
	V := initialConditions["V"]
	F := initialConditions["F"]

	// Extract parameters
	muRef := params["muRef"]
	qpRef := params["qpRef"]
	qsRef := params["qsRef"]
	EaMu := params["EaMu"]
	EaQp := params["EaQp"]
	EaQs := params["EaQs"]
	Kis := params["Kis"]
	Pix := params["Pix"]
	Pmx := params["Pmx"]

	// Storage for results
	results := make([][]float64, 0)

	// Simulation loop
	for t := 0; t < timeSteps; t++ {
		T := temperatureProfile[t] // Get temperature at current time step
		var err error
		X, P, S, V, err = kinethic.KineticModelStep(X, P, S, V, F, T, muRef, qpRef, qsRef, EaMu, EaQp, EaQs, Kis, Pix, Pmx, dt)
		if err != nil {
			fmt.Printf("Error at time step %d: %v\n", t, err)
			break
		}
		results = append(results, []float64{X, P, S, V, T})
	}

	return results
}
