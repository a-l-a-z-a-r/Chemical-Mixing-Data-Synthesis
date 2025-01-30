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
		"X": 0.1,  // Initial biomass (g/L)
		"P": 0.02, // Initial lactic acid (g/L)
		"S": 45.0, // Initial lactose (g/L)
		"V": 0.5,  // Initial volume (L)
		"F": 1.0,  // Flow rate (L/h)
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
		"Pix":   4.8,      // Initial inhibition concentration
		"Pmx":   5.0,      // Maximum inhibitory concentration (set higher than Pix)
	}

	// Temperature profile (oscillating temperature around 300K)
	timeSteps := 100
	var dt float64 = 100.0
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
	defer file.Close() // Ensure the file is closed after writing

	// Create a CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush() // Ensure data is written to file

	// Write the CSV header with spacing for clarity
	header := []string{"Time (s)", "Biomass (g/L)", "Lactic Acid (g/L)", "Lactose (g/L)", "Volume (L)", "Temperature (K)"}
	if err := writer.Write(header); err != nil {
		fmt.Println("Error writing CSV header:", err)
		return
	}

	// Write simulation data to CSV with clear formatting
	for i, res := range results {
		row := []string{
			fmt.Sprintf("%-10.1f", float64(i)*dt), // Align values properly
			fmt.Sprintf("%-12.4f", res[0]),
			fmt.Sprintf("%-15.4f", res[1]),
			fmt.Sprintf("%-15.4f", res[2]),
			fmt.Sprintf("%-10.4f", res[3]),
			fmt.Sprintf("%-10.2f", res[4]),
		}
		if err := writer.Write(row); err != nil {
			fmt.Println("Error writing CSV row:", err)
			return
		}
	}
}
