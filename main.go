package main

import (
	"fmt"
	"math"
)

// KineticModelStep performs one step of the kinetic model simulation
func KineticModelStep(X, P, S, V, F, muMax, alpha, qpMax, qsMax, Kis, Pix, Pmx, dt float64) (float64, float64, float64, float64) {
	// Calculate rate changes
	dX_dt := muMax*X*(1-(P-Pix)/(Pmx-Pix)) + F*X/V
	dP_dt := alpha*dX_dt + qpMax*X*S/(Kis+S) + F*P/V - alpha*X/V
	dS_dt := -qsMax*X*Kis/(Kis+S) + F*S/V
	dV_dt := F
	if V >= 2.0 { // Check if volume exceeds the threshold
		dV_dt = 0.0
	}

	// Update values using Euler's method
	X += dX_dt * dt
	P += dP_dt * dt
	S += dS_dt * dt
	V += dV_dt * dt

	// Prevent negative or invalid values
	if X < 0 {
		X = 0
	}
	if P < 0 {
		P = 0
	}
	if S < 0 {
		S = 0
	}
	if V < 0 {
		V = 0.001 // Prevent division by zero
	}

	return X, P, S, V
}

// CalculatePH calculates pH based on lactic acid concentration using the Henderson-Hasselbalch equation
func CalculatePH(P float64) float64 {
	pKa := 3.86 // pKa value of lactic acid
	Ka := math.Pow(10, -pKa)

	// [A^-] is the concentration of the conjugate base (lactate ion)
	// [HA] is the concentration of the undissociated acid
	// Assuming total lactic acid concentration is the same as P for simplicity
	if P == 0 {
		return 14.0 // Return a high pH if lactic acid concentration is zero
	}
	hPlus := math.Sqrt(Ka * P) // Simplified dissociation equation
	return -math.Log10(hPlus)
}

// SimulateFermentation generates a time series of results from the kinetic model
func SimulateFermentation(muMax, alpha, qpMax, qsMax, Kis, Pix, Pmx, F, dt float64, timeSteps int) [][]float64 {
	// Initial conditions
	X := 0.1  // Biomass concentration
	P := 0.02 // Lactic acid concentration
	S := 45.0 // Lactose concentration
	V := 0.5  // Initial volume

	results := make([][]float64, 0, timeSteps)

	// Perform the simulation
	for t := 0; t < timeSteps; t++ {
		X, P, S, V = KineticModelStep(X, P, S, V, F, muMax, alpha, qpMax, qsMax, Kis, Pix, Pmx, dt)
		pH := CalculatePH(P) // Calculate pH at each step
		results = append(results, []float64{float64(t) * dt, X, P, S, V, pH})
	}

	return results
}

func main() {
	// Simulation settings
	timeSteps := 100
	dt := 0.1

	// Parameters for the fermentation model
	muMax := 1.54e-10
	alpha := 1.33
	qpMax := 3.75e-5
	qsMax := 2.10e-4
	Kis := 5.41e5
	Pix := 4.0
	Pmx := 4.8
	F := 1.0

	// Run simulation for one variation
	results := SimulateFermentation(muMax, alpha, qpMax, qsMax, Kis, Pix, Pmx, F, dt, timeSteps)

	// Print results
	fmt.Println("Time\tX (Biomass)\tP (Lactic Acid)\tS (Lactose)\tV (Volume)\tpH")
	for _, result := range results {
		fmt.Printf("%.2f\t%.6f\t%.6f\t%.6f\t%.6f\t%.2f\n", result[0], result[1], result[2], result[3], result[4], result[5])
	}
}
