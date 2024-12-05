package main

import (
	"reflect"
	"testing"
)

func TestCalculatePH(t *testing.T) {
	type args struct {
		P float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CalculatePH(tt.args.P); got != tt.want {
				t.Errorf("CalculatePH() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKineticModelStep(t *testing.T) {
	type args struct {
		X     float64
		P     float64
		S     float64
		V     float64
		F     float64
		muMax float64
		alpha float64
		qpMax float64
		qsMax float64
		Kis   float64
		Pix   float64
		Pmx   float64
		dt    float64
	}
	tests := []struct {
		name  string
		args  args
		want  float64
		want1 float64
		want2 float64
		want3 float64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2, got3 := KineticModelStep(tt.args.X, tt.args.P, tt.args.S, tt.args.V, tt.args.F, tt.args.muMax, tt.args.alpha, tt.args.qpMax, tt.args.qsMax, tt.args.Kis, tt.args.Pix, tt.args.Pmx, tt.args.dt)
			if got != tt.want {
				t.Errorf("KineticModelStep() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("KineticModelStep() got1 = %v, want %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("KineticModelStep() got2 = %v, want %v", got2, tt.want2)
			}
			if got3 != tt.want3 {
				t.Errorf("KineticModelStep() got3 = %v, want %v", got3, tt.want3)
			}
		})
	}
}

func TestSimulateFermentation(t *testing.T) {
	type args struct {
		muMax     float64
		alpha     float64
		qpMax     float64
		qsMax     float64
		Kis       float64
		Pix       float64
		Pmx       float64
		F         float64
		dt        float64
		timeSteps int
	}
	tests := []struct {
		name string
		args args
		want [][]float64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SimulateFermentation(tt.args.muMax, tt.args.alpha, tt.args.qpMax, tt.args.qsMax, tt.args.Kis, tt.args.Pix, tt.args.Pmx, tt.args.F, tt.args.dt, tt.args.timeSteps); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SimulateFermentation() = %v, want %v", got, tt.want)
			}
		})
	}
}