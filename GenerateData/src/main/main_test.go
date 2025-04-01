package main

import (
	"testing"

	"github.com/simulation"
	"github.com/phcalc"
)

func TestSimulationPipeline(t *testing.T) {
	initialConditions := map[string]float64{
		"X": 0.137,
		"P": 0.024,
		"S": 41.246,
		"V": 500,
		"F": 0.2778,
	}

	params := map[string]float64{
		"muRef": 1.54e-10, "qpRef": 3.75e-5, "qsRef": 2.10e-4,
		"EaMu": 50000.0, "EaQp": 40000.0, "EaQs": 45000.0,
		"Kis": 5.41e5, "Ksp": -27.5, "Inhib": 1.33, "Pix": 4.8, "Pmx": 5.0,
	}

	dt := 1.0
	timeSteps := 10
	tempProfile := make([]float64, timeSteps)
	for i := range tempProfile {
		tempProfile[i] = 300
	}

	results := simulation.SimulateKineticModel(initialConditions, params, tempProfile, timeSteps, dt)

	if len(results) != timeSteps {
		t.Errorf("Expected %d results, got %d", timeSteps, len(results))
	}

	for _, r := range results {
		if r[0] < 0 || r[1] < 0 || r[2] < 0 || r[3] < 0 {
			t.Errorf("Negative values in results: %+v", r)
		}
	}
}
