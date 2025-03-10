package main

import (
	G "gorgonia.org/gorgonia"
	"gorgonia.org/tensor"
)

// LSTM Model Structure
type LSTM struct {
	wix, wih, bias_i *G.Node
	wfx, wfh, bias_f *G.Node
	wox, woh, bias_o *G.Node
	wcx, wch, bias_c *G.Node
	graph            *G.ExprGraph
}

// Create an LSTM Model
func MakeLSTM(g *G.ExprGraph, hiddenSize, inputSize int) LSTM {
	return LSTM{
		wix:    G.NewMatrix(g, tensor.Float64, G.WithShape(hiddenSize, inputSize), G.WithInit(G.GlorotU(1.0)), G.WithName("wix")),
		wih:    G.NewMatrix(g, tensor.Float64, G.WithShape(hiddenSize, hiddenSize), G.WithInit(G.GlorotU(1.0)), G.WithName("wih")),
		bias_i: G.NewVector(g, tensor.Float64, G.WithShape(hiddenSize), G.WithInit(G.Zeroes()), G.WithName("bias_i")),

		wfx:    G.NewMatrix(g, tensor.Float64, G.WithShape(hiddenSize, inputSize), G.WithInit(G.GlorotU(1.0)), G.WithName("wfx")),
		wfh:    G.NewMatrix(g, tensor.Float64, G.WithShape(hiddenSize, hiddenSize), G.WithInit(G.GlorotU(1.0)), G.WithName("wfh")),
		bias_f: G.NewVector(g, tensor.Float64, G.WithShape(hiddenSize), G.WithInit(G.Zeroes()), G.WithName("bias_f")),

		wox:    G.NewMatrix(g, tensor.Float64, G.WithShape(hiddenSize, inputSize), G.WithInit(G.GlorotU(1.0)), G.WithName("wox")),
		woh:    G.NewMatrix(g, tensor.Float64, G.WithShape(hiddenSize, hiddenSize), G.WithInit(G.GlorotU(1.0)), G.WithName("woh")),
		bias_o: G.NewVector(g, tensor.Float64, G.WithShape(hiddenSize), G.WithInit(G.Zeroes()), G.WithName("bias_o")),

		wcx:    G.NewMatrix(g, tensor.Float64, G.WithShape(hiddenSize, inputSize), G.WithInit(G.GlorotU(1.0)), G.WithName("wcx")),
		wch:    G.NewMatrix(g, tensor.Float64, G.WithShape(hiddenSize, hiddenSize), G.WithInit(G.GlorotU(1.0)), G.WithName("wch")),
		bias_c: G.NewVector(g, tensor.Float64, G.WithShape(hiddenSize), G.WithInit(G.Zeroes()), G.WithName("bias_c")),
		graph:  g,
	}
}

func (lstm *LSTM) Activate(input *G.Node) (*G.Node, error) {

	return input, nil
}
