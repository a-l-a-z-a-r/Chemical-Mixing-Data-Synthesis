package main

import (
	"fmt"
	"log"
	"os"

	. "gorgonia.org/gorgonia"
	G "gorgonia.org/gorgonia"
	"gorgonia.org/tensor"
)

/*
	func main() {
		g := NewGraph()
		a := NewScalar(g, Float64, WithName("a"))
		b := NewScalar(g, Float64, WithName("b"))

		c, err := Add(a, b)

		if err != nil {
			log.Fatal(err)
		}

		machine := NewTapeMachine(g)

		Let(a, 1.0)
		Let(b, 2.0)
		if machine.RunAll() != nil {
			log.Fatal(err)
		}

		fmt.Println(c.Value())
		// Output: 3.0”

}

func main() {
	g := NewGraph()
	matB := []float64{0.9, 0.7, 0.4, 0.2}
	matT := tensor.New(tensor.WithBacking(matB), tensor.WithShape(2, 2))
	mat := G.NewMatrix(g, tensor.Float64, G.WithName("W"), G.WithShape(2, 2), G.WithValue(matT))

	vecB := []float64{5, 7}
	vecT := tensor.New(tensor.WithBacking(vecB), tensor.WithShape(2))
	vec := G.NewVector(g, tensor.Float64, G.WithName("x"), G.WithShape(2), G.WithValue(vecT))
	z, err := G.Mul(mat, vec)

	machine := G.NewTapeMachine(g)
	if machine.RunAll() != nil {
		log.Fatal(err)
	}
	fmt.Println(z.Value().Data())
	os.WriteFile("simple_graph.dot", []byte(g.ToDot()), 0644)

	// Output: [9.4 3.4]”
}
*/