package simulation

import (
	"fmt"

	"github.com/kinethic"
	"github.com/phcalc"
)

func SimulateKineticModel(initialConditions map[string]float64, params map[string]float64, temperatureProfile []float64, timeSteps int, dt float64) [][]float64 {
	// Extract initial conditions
	X := initialConditions["X"]
	P := initialConditions["P"] // Lactic acid concentration
	S := initialConditions["S"]
	V := initialConditions["V"]
	F := initialConditions["F"]

	// Extract parameters
	muRef := params["muRef"]
	qpRef := params["qpRef"]
	qsRef := params["qsRef"]
	EaMu := params["EaMu"]
	Ksp := params["Ksp"]
	EaQp := params["EaQp"]
	Inhib := params["Inhib"]
	EaQs := params["EaQs"]
	Kis := params["Kis"]
	Pix := params["Pix"]
	Pmx := params["Pmx"]

	// Buffer capacity (assumed value, should be defined properly)
	bufferCapacity := 10.0 // Adjust based on the actual fermentation conditions

	// Storage for results
	results := make([][]float64, 0)

	// Simulation loop
	for t := 0; t < timeSteps; t++ {
		T := temperatureProfile[t] // Get temperature at current time step
		var err error

		var chn float64
		var chns = &chn
		if V == 0 {
			*chns = 0.0
		} else {
			*chns = 1.0
		}

		X, P, S, V, err = kinethic.KineticModelStep(X, P, S, V, F, T, muRef, qpRef, qsRef, EaMu, EaQp, EaQs, Kis, Pix, Pmx, Inhib, Ksp, *chns, dt)

		if err != nil {
			fmt.Printf("Error at time step %d: %v\n", t, err)
			break
		}

		// Calculate pH using lactic acid concentration
		ph := phcalc.CalculatePH(P, 0, bufferCapacity) // Assuming only lactic acid for now

		// Store results
		results = append(results, []float64{X, P, S, V, T, ph})
	}

	return results
}
