package main

import (
	"fmt"
	G "gorgonia.org/gorgonia"
	"gorgonia.org/tensor"
)

func main() {
	// Load training data
	samples, err := LoadCSV("fermentation_X_0.082.csv")
	if err != nil {
		fmt.Println("Error loading CSV:", err)
		return
	}

	// Train the LSTM model
	TrainLSTM(samples, 100)
	SaveTrainedModel()
}

// TrainLSTM trains the LSTM model using the provided samples and epoch count.
func TrainLSTM(samples []Sample, epochs int) {
	g := G.NewGraph()
	lstm := MakeLSTM(g, 16, 4) // 16 hidden units, 4 inputs

	x := G.NewMatrix(g, tensor.Float64, G.WithShape(len(samples), 4), G.WithName("X"))
	y := G.NewMatrix(g, tensor.Float64, G.WithShape(len(samples), 4), G.WithName("Y"))

	// Compute Predictions
	pred, err := lstm.Activate(x)
	if err != nil {
		fmt.Println("Error in LSTM activation:", err)
		return
	}

	fmt.Println("DEBUG: y shape:", y.Shape())
	fmt.Println("DEBUG: pred shape:", pred.Shape())

	// Ensure y and pred have matching shapes before subtracting
	if !y.Shape().Eq(pred.Shape()) {
		fmt.Println("Error: y and pred have mismatched shapes:", y.Shape(), pred.Shape())
		return
	}

	// Compute Loss Function (MSE)
	diff, err := G.Sub(y, pred)
	if err != nil {
		fmt.Println("Error computing difference:", err)
		return
	}

	squared, err := G.Square(diff)
	if err != nil {
		fmt.Println("Error computing squared loss:", err)
		return
	}

	mean, err := G.Mean(squared)
	if err != nil {
		fmt.Println("Error computing mean loss:", err)
		return
	}

	loss := mean

	// Compute Gradients
	_, err = G.Grad(loss, lstm.wix, lstm.wih, lstm.bias_i, lstm.wfx, lstm.wfh, lstm.bias_f)
	if err != nil {
		panic(err)
	}

	// Create TapeMachine
	vm := G.NewTapeMachine(g, G.BindDualValues(lstm.wix, lstm.wih, lstm.bias_i, lstm.wfx, lstm.wfh, lstm.bias_f))
	solver := G.NewAdamSolver(G.WithLearnRate(0.01))

	// Training Loop
	for i := 0; i < epochs; i++ {
		if err := vm.RunAll(); err != nil {
			panic(err)
		}
		valueGrads := []G.ValueGrad{
			lstm.wix, lstm.wih, lstm.bias_i,
			lstm.wfx, lstm.wfh, lstm.bias_f,
		}
		solver.Step(valueGrads)
		vm.Reset()
		fmt.Printf("Epoch %d: Loss = %v\n", i+1, loss.Value())
	}
}

// SaveTrainedModel is a stub function to save the model.
func SaveTrainedModel() {
	fmt.Println("Saving trained model (implement logic)")
}
