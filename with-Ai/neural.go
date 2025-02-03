package main

import (
	"fmt"
	"log"
	"math"

	G "gorgonia.org/gorgonia"
	"gorgonia.org/tensor"
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
	R := 8.314
	TRef := 298.15

	if temp <= 0 {
		return fmt.Errorf("temperature T must be greater than 0; got %f", temp)
	}

	// Temperature-dependent parameters
	muMax := muRef * math.Exp(-EaMu/R*(1/temp-1/TRef))
	qpMax := qpRef * math.Exp(-EaQp/R*(1/temp-1/TRef))
	qsMax := qsRef * math.Exp(-EaQs/R*(1/temp-1/TRef))

	// Clamp rates
	muMax = math.Max(muMax, 1e-6)
	qpMax = math.Max(qpMax, 1e-6)
	qsMax = math.Max(qsMax, 1e-6)

	dX_dt := muMax*fs.X*(1-(fs.P-fs.Pix)/(fs.Pmx-fs.Pix)) + fs.F*fs.X/fs.V
	dP_dt := fs.alpha*dX_dt + qpMax*fs.X*fs.S/(fs.Kis+fs.S) + fs.F*fs.P/fs.V - fs.alpha*fs.X/fs.V
	dS_dt := -qsMax*fs.X*fs.Kis/(fs.Kis+fs.S) + fs.F*fs.S/fs.V
	dV_dt := fs.F
	if fs.V >= 2.0 {
		dV_dt = 0.0
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
	pKa := 3.86
	Ka := math.Pow(10, -pKa)

	if fs.P == 0 {
		return 14.0
	}
	hPlus := math.Sqrt(Ka * fs.P)
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
		})
	}
}

// NeuralNet defines a simple feedforward neural network
type NeuralNet struct {
	g          *G.ExprGraph
	w0, b0     *G.Node
	w1, b1     *G.Node
	prediction *G.Node
	loss       *G.Node
}

// NewNeuralNet initializes a simple feedforward neural network
func NewNeuralNet(g *G.ExprGraph, inputSize, hiddenSize, outputSize int) *NeuralNet {
	w0 := G.NewMatrix(g, tensor.Float64, G.WithShape(inputSize, hiddenSize), G.WithName("w0"), G.WithInit(G.GlorotN(1)))
	b0 := G.NewVector(g, tensor.Float64, G.WithShape(hiddenSize), G.WithName("b0"), G.WithInit(G.Zeroes()))
	w1 := G.NewMatrix(g, tensor.Float64, G.WithShape(hiddenSize, outputSize), G.WithName("w1"), G.WithInit(G.GlorotN(1)))
	b1 := G.NewVector(g, tensor.Float64, G.WithShape(outputSize), G.WithName("b1"), G.WithInit(G.Zeroes()))
	return &NeuralNet{g: g, w0: w0, b0: b0, w1: w1, b1: b1}
}

// Forward performs a forward pass of the neural network
func (nn *NeuralNet) Forward(x *G.Node) *G.Node {
	// First layer: X * W0 + B0 (bias broadcasted)
	biasBroadcasted := G.Must(G.BroadcastAdd(G.Must(G.Mul(x, nn.w0)), nn.b0, nil, []byte{0}))
	hidden := G.Must(G.Rectify(biasBroadcasted))

	// Second layer: Hidden * W1 + B1 (bias broadcasted)
	outputBiasBroadcasted := G.Must(G.BroadcastAdd(G.Must(G.Mul(hidden, nn.w1)), nn.b1, nil, []byte{0}))
	nn.prediction = outputBiasBroadcasted

	return nn.prediction
}

// SetLoss sets the loss function (mean squared error)
func (nn *NeuralNet) SetLoss(y *G.Node) {
	nn.loss = G.Must(G.Mean(G.Must(G.Square(G.Must(G.Sub(nn.prediction, y))))))
}
func NewSimulation(muMax, alpha, qpMax, qsMax, Kis, Pix, Pmx, F, dt float64, timeSteps int) *Simulation {
	system := FermentationSystem{
		X:     0.1,
		P:     0.02,
		S:     45.0,
		V:     0.5,
		muMax: muMax,
		alpha: alpha,
		qpMax: qpMax,
		qsMax: qsMax,
		Kis:   Kis,
		Pix:   Pix,
		Pmx:   Pmx,
		F:     F,
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
	timeSteps := 10000
	dt := 0.1
	temperatureProfile := make([]float64, timeSteps)
	for i := 0; i < timeSteps; i++ {
		temperatureProfile[i] = 300 + 5*math.Sin(2*math.Pi*float64(i)/float64(timeSteps))
	}

	// Fermentation parameters
	muRef := 1.54e-10
	alpha := 1.33
	qpRef := 3.75e-5
	qsRef := 2.10e-4
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

	// Prepare training data
	inputs := []float64{}
	outputs := []float64{}
	for _, result := range simulation.Results {
		inputs = append(inputs, result[6], result[3], result[2]) // Temp, Lactose, Lactic Acid
		outputs = append(outputs, result[5])                     // pH
	}

	xTensor := tensor.New(tensor.WithShape(timeSteps, 3), tensor.WithBacking(inputs))
	yTensor := tensor.New(tensor.WithShape(timeSteps, 1), tensor.WithBacking(outputs))

	// Create computation graph
	g := G.NewGraph()

	// Placeholders for input and output
	x := G.NewMatrix(g, tensor.Float64, G.WithShape(timeSteps, 3), G.WithName("x"))
	y := G.NewMatrix(g, tensor.Float64, G.WithShape(timeSteps, 1), G.WithName("y"))

	// Create the neural network
	nn := NewNeuralNet(g, 3, 10, 1)
	nn.Forward(x)
	nn.SetLoss(y)

	// Compute gradients
	if _, err := G.Grad(nn.loss, nn.w0, nn.b0, nn.w1, nn.b1); err != nil {
		log.Fatalf("Failed to compute gradients: %v", err)
	}

	// Training loop
	vm := G.NewTapeMachine(g, G.BindDualValues(nn.w0, nn.b0, nn.w1, nn.b1))
	solver := G.NewVanillaSolver(G.WithLearnRate(0.01))
	for epoch := 0; epoch < 500; epoch++ {
		G.Let(x, xTensor)
		G.Let(y, yTensor)

		if err := vm.RunAll(); err != nil {
			log.Fatalf("Failed to run VM: %v", err)
		}
		solver.Step([]G.ValueGrad{nn.w0, nn.b0, nn.w1, nn.b1})

		vm.Reset()

		if epoch%50 == 0 {
			fmt.Printf("Epoch %d: Loss: %v\n", epoch, nn.loss.Value())
		}
	}

	fmt.Println("Training complete.")
}
