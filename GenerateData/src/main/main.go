package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"os"

	"github.com/simulation"
)

// Function to run the simulation for a given set of initial conditions
func runSimulation(initialConditions map[string]float64, params map[string]float64, temperatureProfile []float64, timeSteps int, dt float64, fileName string) {
	// Run the simulation
	results := simulation.SimulateKineticModel(initialConditions, params, temperatureProfile, timeSteps, dt)

	// Create a CSV file
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatalf("Error creating CSV file %s: %v", fileName, err)
	}
	defer file.Close()

	// Create a CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write the CSV header
	header := []string{"Time (s)", "Biomass (g/mL)", "Lactic Acid (g/mL)", "Lactose (g/mL)", "Volume (mL)", "Temperature (K)", "pH"}
	if err := writer.Write(header); err != nil {
		log.Fatalf("Error writing CSV header for %s: %v", fileName, err)
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
			log.Fatalf("Error writing CSV row for %s: %v", fileName, err)
		}
	}

	fmt.Printf("Simulation completed. Results saved in '%s'.\n", fileName)
}

func main() {
	// List of initial conditions for different biomass values
	initialConditionsList := []map[string]float64{
		{"X": 0.137, "P": 41.246, "S": 0.024, "V": 500, "F": 0.2778},
		{"X": 0.040, "P": 44.102, "S": 0.023, "V": 1500, "F": 0.2778},
		{"X": 0.082, "P": 45.563, "S": 0.021, "V": 1500, "F": 0.2778},
	}

	// Parameters for simulation
	params := map[string]float64{
		"muRef": 1.54e-10, "qpRef": 3.75e-5, "qsRef": 2.10e-4,
		"EaMu": 50000.0, "EaQp": 40000.0, "EaQs": 45000.0,
		"Kis": 5.41e5, "Ksp": -27.50, "Inhib": 1.33,
		"Pix": 4.8, "Pmx": 5.0,
	}

	// Generate a temperature profile
	timeSteps := 18000
	dt := 1.0
	temperatureProfile := make([]float64, timeSteps)
	for i := 0; i < timeSteps; i++ {
		temperatureProfile[i] = 300 + 5*math.Sin(2*math.Pi*float64(i)/float64(timeSteps))
	}

	// Run simulation for each set of initial conditions
	for _, initialConditions := range initialConditionsList {
		// Create a unique filename based on initial biomass value
		fileName := fmt.Sprintf("fermentation_X_%.3f.csv", initialConditions["X"])

		// Run the simulation and save results
		runSimulation(initialConditions, params, temperatureProfile, timeSteps, dt, fileName)
	}
}
