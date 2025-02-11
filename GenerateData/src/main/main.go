package main

import (
	"encoding/csv"
	"fmt"
	"math"
	"os"

	"github.com/simulation"
)

func main() {
	// Initial conditions
	initialConditions := map[string]float64{
		"X": 0.137,  // Initial biomass (g/L)
		"P": 0.024,  // Initial lactic acid (g/L)
		"S": 41.246, // Initial lactose (g/L)
		"V": 500,    // Initial volume (mL)
		"F": 0.2778, // Flow rate (L/h)
	}

	// Parameters
	params := map[string]float64{
		"muRef": 1.54e-10, // Reference growth rate
		"qpRef": 3.75e-5,  // Reference lactic acid production rate
		"qsRef": 2.10e-4,  // Reference lactose utilization rate
		"EaMu":  50000.0,  // Activation energy for growth rate (J/mol)
		"EaQp":  40000.0,  // Activation energy for lactic acid production (J/mol)
		"EaQs":  45000.0,  // Activation energy for lactose utilization (J/mol)
		"Kis":   5.41e5,   // Lactose limitation constant
		"Ksp":   -27.50,   // Unknown equilibrium constant
		"Inhib": 1.33,     // Inhibition constant
		"Pix":   4.8,      // Initial inhibition concentration
		"Pmx":   5.0,      // Maximum inhibitory concentration
	}

	// Temperature profile (oscillating temperature around 300K)
	timeSteps := 5399
	var dt float64 = 1 // Measurements taken every hour for 25 hours
	temperatureProfile := make([]float64, timeSteps)
	for i := 0; i < timeSteps; i++ {
		temperatureProfile[i] = 300 + 5*math.Sin(2*math.Pi*float64(i)/float64(timeSteps)) // Oscillating temp
	}

	// Run simulation
	results := simulation.SimulateKineticModel(initialConditions, params, temperatureProfile, timeSteps, dt)

	// Create a CSV file
	file, err := os.Create("fermentation_results.csv")
	if err != nil {
		fmt.Println("Error creating CSV file:", err)
		return
	}
	defer file.Close()

	// Create a CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write the CSV header
	header := []string{"Time (s)", "Biomass (g/mL)", "Lactic Acid (g/mL)", "Lactose (g/mL)", "Volume (mL)", "Temperature (K)", "pH"}
	if err := writer.Write(header); err != nil {
		fmt.Println("Error writing CSV header:", err)
		return
	}

	// Write simulation data to CSV
	for i, res := range results {
		row := []string{
			fmt.Sprintf("%.1f", float64(i)*dt), // Time
			fmt.Sprintf("%.4f", res[0]),        // Biomass
			fmt.Sprintf("%.4f", res[1]),        // Lactic Acid
			fmt.Sprintf("%.4f", res[2]),        // Lactose
			fmt.Sprintf("%.4f", res[3]),        // Volume
			fmt.Sprintf("%.2f", res[4]),        // Temperature
			fmt.Sprintf("%.2f", res[5]),        // pH
		}
		if err := writer.Write(row); err != nil {
			fmt.Println("Error writing CSV row:", err)
			return
		}
	}

	fmt.Println("Simulation completed. Results saved in 'fermentation_results.csv'.")
}
