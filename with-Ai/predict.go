package main

import (
	"fmt"
	G "gorgonia.org/gorgonia"
	"gorgonia.org/tensor"
)

// Predict using trained LSTM
func PredictLSTM(lstm *LSTM, input []float64) float64 {
	tensorInput := tensor.New(tensor.WithShape(1, 4), tensor.Of(tensor.Float64), tensor.WithBacking(input))
	output, _ := lstm.Activate(G.NewMatrix(lstm.graph, tensor.Float64, G.WithShape(1, 4), G.WithValue(tensorInput)))

	return output.Value().Data().([]float64)[0]
}

// Load Model
func LoadTrainedModel() *LSTM {
	fmt.Println("Loading trained model (implement logic)")
	return &LSTM{}
}
