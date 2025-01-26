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
func (km *KineticModel) Step(fs *FermentationSystem, temp float64, muRef, qpRef, qsRef, EaMu, EaQp, EaQs float64) error {
	// Constants
	R := 8.314
	TRef := 298.15

	// Ensure temperature is valid
	if temp <= 0 {
		return fmt.Errorf("temperature T must be greater than 0; got %f", temp)
	}

	// Calculate temperature-dependent rates
	muMax := muRef * math.Exp(-EaMu/R*(1/temp-1/TRef))
	qpMax := qpRef * math.Exp(-EaQp/R*(1/temp-1/TRef))
	qsMax := qsRef * math.Exp(-EaQs/R*(1/temp-1/TRef))

	// Clamp rates
	muMax = math.Max(muMax, 1e-6)
	qpMax = math.Max(qpMax, 1e-6)
	qsMax = math.Max(qsMax, 1e-6)

	// Existing kinetic equations
	dX_dt := muMax*fs.X*(1-(fs.P-fs.Pix)/(fs.Pmx-fs.Pix)) + fs.F*fs.X/fs.V
	dP_dt := fs.alpha*dX_dt + qpMax*fs.X*fs.S/(fs.Kis+fs.S) + fs.F*fs.P/fs.V - fs.alpha*fs.X/fs.V
	dS_dt := -qsMax*fs.X*fs.Kis/(fs.Kis+fs.S) + fs.F*fs.S/fs.V
	dV_dt := fs.F
	if fs.V >= 2.0 { // If volume is greater than 2.0, stop adding more
		fs.F = 0.0 // we had dV_dt = 0.0 in another version of the code, because of this change P becomes 0.0 in the results, we have to think about that,X and S got o zero too if we change below
		//fs.X = 0.0	//" which is the same for X, P, S"
		//fs.S = 0.0
	}

	// Euler integration
	fs.X += dX_dt * fs.dt
	fs.P += dP_dt * fs.dt
	fs.S += dS_dt * fs.dt
	fs.V += dV_dt * fs.dt

	// Prevent invalid values
	fs.X = math.Max(fs.X, 0)
	fs.P = math.Max(fs.P, 0)
	fs.S = math.Max(fs.S, 0)
	fs.V = math.Max(fs.V, 0.001)

	return nil
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
	System             FermentationSystem
	KineticModel       KineticModel
	PHCalculator       PHCalculator
	TemperatureProfile []float64
	TimeSteps          int
	Results            [][]float64
}

// Run executes the simulation over the specified time steps
func (sim *Simulation) Run(muRef, qpRef, qsRef, EaMu, EaQp, EaQs float64) {
	for t := 0; t < sim.TimeSteps; t++ {
		temp := sim.TemperatureProfile[t]
		err := sim.KineticModel.Step(&sim.System, temp, muRef, qpRef, qsRef, EaMu, EaQp, EaQs)
		if err != nil {
			fmt.Printf("Error at step %d: %v\n", t, err)
			break
		}
		pH := sim.PHCalculator.Calculate(&sim.System)
		sim.Results = append(sim.Results, []float64{
			float64(t) * sim.System.dt,
			sim.System.X,
			sim.System.P,
			sim.System.S,
			sim.System.V,
			pH,
			temp,
			//sim.System.F,
		})
	}
}

// NewSimulation initializes and returns a new Simulation instance
func NewSimulation(muMax, alpha, qpMax, qsMax, Kis, Pix, Pmx, F, dt float64, timeSteps int) *Simulation {
	system := FermentationSystem{
		X:     0.1,   // Biomass concentration
		P:     0.02,  // Lactic acid concentration
		S:     45.0,  // Lactose concentration
		V:     0.5,   // Initial volume
		muMax: muMax, // Maximum growth rate
		alpha: alpha, // Inhibition constant
		qpMax: qpMax, // Maximum specific lactic acid production rate
		qsMax: qsMax, // Maximum specific lactose utilisation rate
		Kis:   Kis,
		Pix:   Pix, // Threshold lactate concentration before any inhibition occurs
		Pmx:   Pmx, // Maximum inhibitory value
		F:     F,   // Flow-in rate
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
	// Simulation parameters
	timeSteps := 100
	dt := 0.1

	// We have to variate the 300 in the loop below
	// Temperature profile
	temperatureProfile := make([]float64, timeSteps)
	for i := 0; i < timeSteps; i++ {
		temperatureProfile[i] = 300 + 5*math.Sin(2*math.Pi*float64(i)/float64(timeSteps)) // Oscillating temp
	}

	// Fermentation parameters - from literature: "pH prediction for a semi-batch cream cheese
	// fermentation using a grey-box model", table 2

	muRef := 1.54e-10 // refers to muMax
	alpha := 1.33
	qpRef := 3.75e-5 // refers to qpMax
	qsRef := 2.10e-4 // refers to qsMax
	EaMu := 50000.0
	EaQp := 40000.0
	EaQs := 45000.0
	Kis := 5.41e5
	Pix := 4.0
	Pmx := 4.8
	F := 1.0

	// Initialize simulation
	simulation := NewSimulation(muRef, alpha, qpRef, qsRef, Kis, Pix, Pmx, F, dt, timeSteps)
	simulation.TemperatureProfile = temperatureProfile
	simulation.Run(muRef, qpRef, qsRef, EaMu, EaQp, EaQs)

	// Print results
	fmt.Println("Time\tX\tP\tS\tV\tpH\tTemp")
	for _, result := range simulation.Results {
		fmt.Printf("%.2f\t%.6f\t%.6f\t%.6f\t%.6f\t%.2f\t%.2f\n", result[0], result[1], result[2], result[3], result[4], result[5], result[6])
	}
}
