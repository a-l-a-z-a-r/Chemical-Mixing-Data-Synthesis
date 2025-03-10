package main

import (
	"fmt"
)

func main() {
	var mode string
	fmt.Println("Choose mode: (train/predict)")
	fmt.Scanln(&mode)

	if mode == "train" {
		samples, _ := LoadCSV("fermentation_X_0.082.csv")
		TrainLSTM(samples, 100)
		SaveTrainedModel()
	} else if mode == "predict" {
		lstm := LoadTrainedModel()
		input := []float64{0.2, 0.03, 0.6, 1.5}
		pH := PredictLSTM(lstm, input)
		fmt.Printf("Predicted pH: %.2f\n", pH)
	} else {
		fmt.Println("Invalid mode! Use 'train' or 'predict'.")
	}
}
