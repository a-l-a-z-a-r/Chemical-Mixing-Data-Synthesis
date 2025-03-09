package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"

	G "gorgonia.org/gorgonia"
	"gorgonia.org/tensor"
)

// Hyperparameters
const (
	inputSize    = 4  // X, P, S, V
	hiddenSize   = 16 // LSTM hidden units
	outputSize   = 1  // pH prediction
	epochs       = 500
	learningRate = 0.001
	modelFile    = "lstm_model.json"
)

var mu sync.Mutex // Ensure thread safety for concurrent API calls

// LSTM Model
type LSTM struct {
	g      *G.ExprGraph
	w0, w1 *G.Node
	b0, b1 *G.Node
	pred   *G.Node
	vm     G.VM
}

// NewLSTM initializes an LSTM model
func NewLSTM() *LSTM {
	g := G.NewGraph()

	// Weights and biases
	w0 := G.NewMatrix(g, tensor.Float32, G.WithShape(inputSize, hiddenSize), G.WithInit(G.GlorotN(1.0)))
	b0 := G.NewMatrix(g, tensor.Float32, G.WithShape(1, hiddenSize), G.WithInit(G.Zeroes()))
	w1 := G.NewMatrix(g, tensor.Float32, G.WithShape(hiddenSize, outputSize), G.WithInit(G.GlorotN(1.0)))
	b1 := G.NewMatrix(g, tensor.Float32, G.WithShape(1, outputSize), G.WithInit(G.Zeroes()))

	return &LSTM{g: g, w0: w0, b0: b0, w1: w1, b1: b1}
}
func (l *LSTM) Forward(x *G.Node) (*G.Node, error) {
	batchSize := x.Shape()[0]

	biasReshaped, err := G.Reshape(l.b0, tensor.Shape{batchSize, hiddenSize})
	if err != nil {
		return nil, err
	}

	hid := G.Must(G.Add(G.Must(G.Mul(x, l.w0)), biasReshaped))
	hid = G.Must(G.Tanh(hid))

	out := G.Must(G.Add(G.Must(G.Mul(hid, l.w1)), l.b1))
	return out, nil
}

// Train LSTM model
func (l *LSTM) Train(X, Y *tensor.Dense) {
	mu.Lock()
	defer mu.Unlock()

	g := l.g
	x := G.NewMatrix(g, tensor.Float32, G.WithShape(X.Shape()...), G.WithValue(X))
	y := G.NewMatrix(g, tensor.Float32, G.WithShape(Y.Shape()...), G.WithValue(Y))

	pred, err := l.Forward(x)
	if err != nil {
		log.Fatal(err)
	}

	loss := G.Must(G.Mean(G.Must(G.Square(G.Must(G.Sub(y, pred))))))

	// Compute gradients
	_, err = G.Grad(loss, l.w0, l.b0, l.w1, l.b1)
	if err != nil {
		log.Fatal(err)
	}

	// Training function
	vm := G.NewTapeMachine(g, G.BindDualValues(l.w0, l.b0, l.w1, l.b1))
	solver := G.NewAdamSolver(G.WithLearnRate(learningRate))

	// Training loop
	for i := 0; i < epochs; i++ {
		if err := vm.RunAll(); err != nil {
			log.Fatal(err)
		}

		// **Correct gradient update**
		solver.Step([]G.ValueGrad{l.w0, l.b0, l.w1, l.b1})

		vm.Reset()
		fmt.Printf("Epoch %d: Loss = %v\n", i+1, loss.Value())
	}

	// Save model after training
	l.SaveModel()
}

// Predict pH
func (l *LSTM) Predict(X *tensor.Dense) float64 {
	mu.Lock()
	defer mu.Unlock()

	g := l.g
	x := G.NewMatrix(g, tensor.Float32, G.WithShape(X.Shape()...), G.WithValue(X))

	pred, err := l.Forward(x)
	if err != nil {
		log.Fatal(err)
	}

	vm := G.NewTapeMachine(g)
	if err := vm.RunAll(); err != nil {
		log.Fatal(err)
	}

	result := pred.Value().Data().([]float32)
	return float64(result[0])
}

// Save Model
func (l *LSTM) SaveModel() {
	file, _ := os.Create(modelFile)
	defer file.Close()

	data := map[string]interface{}{
		"w0": l.w0.Value().Data(),
		"b0": l.b0.Value().Data(),
		"w1": l.w1.Value().Data(),
		"b1": l.b1.Value().Data(),
	}

	json.NewEncoder(file).Encode(data)
	fmt.Println("Model saved successfully!")
}

// Load Model
func (l *LSTM) LoadModel() {
	file, err := os.Open(modelFile)
	if err != nil {
		fmt.Println("No previous model found. Training from scratch.")
		return
	}
	defer file.Close()

	data := make(map[string]interface{})
	json.NewDecoder(file).Decode(&data)

	l.w0 = G.NewMatrix(l.g, tensor.Float32, G.WithShape(inputSize, hiddenSize), G.WithValue(data["w0"]))
	l.b0 = G.NewMatrix(l.g, tensor.Float32, G.WithShape(1, hiddenSize), G.WithValue(data["b0"]))
	l.w1 = G.NewMatrix(l.g, tensor.Float32, G.WithShape(hiddenSize, outputSize), G.WithValue(data["w1"]))
	l.b1 = G.NewMatrix(l.g, tensor.Float32, G.WithShape(1, outputSize), G.WithValue(data["b1"]))

	fmt.Println("Model loaded successfully!")
}

func main() {
	lstm := NewLSTM()
	lstm.LoadModel()

	// Load data
	file, err := os.Open("fermentation_X_0.082.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	lines, _ := reader.ReadAll()

	Xdata := make([]float32, 0)
	Ydata := make([]float32, 0)

	for _, line := range lines[1:] {
		x, _ := strconv.ParseFloat(line[1], 32)
		p, _ := strconv.ParseFloat(line[2], 32)
		s, _ := strconv.ParseFloat(line[3], 32)
		v, _ := strconv.ParseFloat(line[4], 32)
		ph, _ := strconv.ParseFloat(line[5], 32)

		Xdata = append(Xdata, float32(x), float32(p), float32(s), float32(v))
		Ydata = append(Ydata, float32(ph))
	}

	X := tensor.New(tensor.WithShape(len(Xdata)/4, 4), tensor.Of(tensor.Float32), tensor.WithBacking(Xdata))
	Y := tensor.New(tensor.WithShape(len(Ydata), 1), tensor.Of(tensor.Float32), tensor.WithBacking(Ydata))

	lstm.Train(X, Y)

	// Example prediction
	sampleX := tensor.New(tensor.WithShape(1, 4), tensor.Of(tensor.Float32), tensor.WithBacking([]float32{0.5, 0.3, 0.2, 1.0}))
	predictedPH := lstm.Predict(sampleX)
	fmt.Printf("Predicted pH: %.4f\n", predictedPH)
}
