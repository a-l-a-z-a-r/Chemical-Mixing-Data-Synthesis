package kinethic

import (
	"fmt"
	"math"

	"golang.org/x/tools/go/analysis/passes/nilfunc"
)

// Hello returns a greeting for the named person.
func KineticModelStep(X, P, S, V, F, T, muRef, qpRef, qsRef, EaMu, EaQp, EaQs, Kis, Pix, Pmx, dt float64) (float64, float64, float64, float64, error) {
	// Constants
	R := 8.314     // Gas constant (J/mol/K)
	TRef := 298.15 // Reference temperature (K)

	if T <= 0 {
		return X, P, S, V, fmt.Errorf("temperature T must be greater than 0; got %f", T)
	}

	if F >= 2 {
		return X, P, S, V, fmt.Errorf("Flow rate can't be larger than 2 Liters  %f", F)
	}
	// Temperature-dependent parameters using Arrhenius equation
	muMax := muRef * math.Exp(-EaMu/R*(1/T-1/TRef))
	qpMax := qpRef * math.Exp(-EaQp/R*(1/T-1/TRef))
	qsMax := qsRef * math.Exp(-EaQs/R*(1/T-1/TRef))

	// Differential equations
	dX := muMax*X*(1-(P-Pix)/(Pmx-Pix)) + F*X/V
	dP := qpMax*dX + qpMax*X*S/(Kis+S) + F*P/V
	dS := -qsMax*X*Kis/(Kis+S) + F*S/V
	dV := F

	// Update variables using Euler's method
	X += dX * dt
	P += dP * dt
	S += dS * dt
	V += dV * dt

	return X, P, S, V, nil
}
