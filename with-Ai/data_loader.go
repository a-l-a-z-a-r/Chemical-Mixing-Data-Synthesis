package main

import (
	"encoding/csv"

	"os"
	"strconv"
)

// Structure for the dataset
type Sample struct {
	X, P, S, V, PH float64
}

// Load CSV Data
func LoadCSV(filename string) ([]Sample, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var samples []Sample
	for i, line := range lines {
		if i == 0 {
			continue // Skip header
		}
		x, _ := strconv.ParseFloat(line[0], 64)
		p, _ := strconv.ParseFloat(line[1], 64)
		s, _ := strconv.ParseFloat(line[2], 64)
		v, _ := strconv.ParseFloat(line[3], 64)
		ph, _ := strconv.ParseFloat(line[4], 64)

		samples = append(samples, Sample{X: x, P: p, S: s, V: v, PH: ph})
	}

	return samples, nil
}
