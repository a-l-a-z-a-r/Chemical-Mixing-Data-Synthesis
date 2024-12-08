package main

import (
	"fmt"
	"math"
)

// FermentationSystem represents the kinetic model and its state
type FermentationSystem struct {
	X     float64 // Biomass concentration
	P     float64 // Lactic acid concentration
	S     float64 // Lactose concentration
	V     float64 // Volume
	muMax float64
	alpha float64
	qpMax float64
	qsMax float64
	Kis   float64
	Pix   float64
	Pmx   float64
	F     float64
	dt    float64
}

// KineticModel encapsulates the logic for a single kinetic step
type KineticModel struct{}

// Step performs one step of the kinetic model simulation
func (km *KineticModel) Step(fs *FermentationSystem) {
	// Calculate rate changes
	dX_dt := fs.muMax*fs.X*(1-(fs.P-fs.Pix)/(fs.Pmx-fs.Pix)) + fs.F*fs.X/fs.V
	dP_dt := fs.alpha*dX_dt + fs.qpMax*fs.X*fs.S/(fs.Kis+fs.S) + fs.F*fs.P/fs.V - fs.alpha*fs.X/fs.V
	dS_dt := -fs.qsMax*fs.X*fs.Kis/(fs.Kis+fs.S) + fs.F*fs.S/fs.V
	dV_dt := fs.F
	if fs.V >= 2.0 { // Check if volume exceeds the threshold
		dV_dt = 0.0
	}

	// Update values using Euler's method
	fs.X += dX_dt * fs.dt
	fs.P += dP_dt * fs.dt
	fs.S += dS_dt * fs.dt
	fs.V += dV_dt * fs.dt

	// Prevent negative or invalid values
	if fs.X < 0 {
		fs.X = 0
	}
	if fs.P < 0 {
		fs.P = 0
	}
	if fs.S < 0 {
		fs.S = 0
	}
	if fs.V < 0 {
		fs.V = 0.001 // Prevent division by zero
	}
}

// PHCalculator encapsulates the logic for calculating pH
type PHCalculator struct{}

// Calculate calculates pH based on lactic acid concentration
func (pc *PHCalculator) Calculate(fs *FermentationSystem) float64 {
	pKa := 3.86 // pKa value of lactic acid
	Ka := math.Pow(10, -pKa)

	if fs.P == 0 {
		return 14.0 // Return a high pH if lactic acid concentration is zero
	}
	hPlus := math.Sqrt(Ka * fs.P) // Simplified dissociation equation
	return -math.Log10(hPlus)
}

// Simulation represents the entire fermentation process
type Simulation struct {
	System       FermentationSystem
	KineticModel KineticModel
	PHCalculator PHCalculator
	TimeSteps    int
	Results      [][]float64
}

// Run executes the simulation over the specified time steps
func (sim *Simulation) Run() {
	for t := 0; t < sim.TimeSteps; t++ {
		sim.KineticModel.Step(&sim.System)
		pH := sim.PHCalculator.Calculate(&sim.System)
		sim.Results = append(sim.Results, []float64{
			float64(t) * sim.System.dt,
			sim.System.X,
			sim.System.P,
			sim.System.S,
			sim.System.V,
			pH,
		})
	}
}

// NewSimulation initializes and returns a new Simulation instance
func NewSimulation(muMax, alpha, qpMax, qsMax, Kis, Pix, Pmx, F, dt float64, timeSteps int) *Simulation {
	system := FermentationSystem{
		X:     0.1,  // Biomass concentration
		P:     0.02, // Lactic acid concentration
		S:     45.0, // Lactose concentration
		V:     0.5,  // Initial volume
		muMax: muMax,
		alpha: alpha,
		qpMax: qpMax,
		qsMax: qsMax,
		Kis:   Kis,
		Pix:   Pix,
		Pmx:   Pmx,
		F:     F,
		dt:    dt,
	}

	return &Simulation{
		System:       system,
		KineticModel: KineticModel{},
		PHCalculator: PHCalculator{},
		TimeSteps:    timeSteps,
		Results:      make([][]float64, 0, timeSteps),
	}
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

	// Initialize and run the simulation
	simulation := NewSimulation(muMax, alpha, qpMax, qsMax, Kis, Pix, Pmx, F, dt, timeSteps)
	simulation.Run()

	// Print results
	fmt.Println("Time\tX (Biomass)\tP (Lactic Acid)\tS (Lactose)\tV (Volume)\tpH")
	for _, result := range simulation.Results {
		fmt.Printf("%.2f\t%.6f\t%.6f\t%.6f\t%.6f\t%.2f\n", result[0], result[1], result[2], result[3], result[4], result[5])
	}
}
