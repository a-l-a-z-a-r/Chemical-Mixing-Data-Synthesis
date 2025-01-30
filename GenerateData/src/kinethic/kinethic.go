package kinethic

import (
	"fmt"
	"math"
)

// Hello returns a greeting for the named person.
func KineticModelStep(X, P, S, V, F, T, muRef, qpRef, qsRef, EaMu, EaQp, EaQs, Kis, Pix, Pmx, dt float64) (float64, float64, float64, float64, error) {
	// Constants
	R := 8.314     // Gas constant (J/mol/K)
	TRef := 298.15 // Reference temperature (K)

	if T <= 0 {
		return X, P, S, V, fmt.Errorf("temperature T must be greater than 0; got %f", T)
	}

	// Temperature-dependent parameters using Arrhenius equation
	muMax := muRef * math.Exp(-EaMu/R*(1/T-1/TRef))
	qpMax := qpRef * math.Exp(-EaQp/R*(1/T-1/TRef))
	qsMax := qsRef * math.Exp(-EaQs/R*(1/T-1/TRef))

	// Clamp very small or invalid rates
	if muMax < 1e-6 {
		muMax = 1e-6
	}
	if qpMax < 1e-6 {
		qpMax = 1e-6
	}
	if qsMax < 1e-6 {
		qsMax = 1e-6
	}

	// Check for valid inhibition parameters
	if math.Abs(Pmx-Pix) < 1e-6 {
		return X, P, S, V, fmt.Errorf("Pmx and Pix are too close or identical: Pmx=%f, Pix=%f", Pmx, Pix)
	}

	// Differential equations
	dX := muMax*X*(1-(P-Pix)/(Pmx-Pix)) + F*X/V
	dP := qpMax*dX + qpMax*X*S/(Kis+S) + F*P/V
	dS := -qsMax*X*S/(Kis+S) + F*S/V
	dV := F

	// Clamp extreme rate changes
	if dX > 1e6 {
		dX = 1e6
	}
	if dP > 1e6 {
		dP = 1e6
	}

	// Check for invalid or infinite rates
	if math.IsNaN(dX) || math.IsNaN(dP) || math.IsNaN(dS) || math.IsInf(dX, 0) || math.IsInf(dP, 0) || math.IsInf(dS, 0) {
		return X, P, S, V, fmt.Errorf("invalid rate changes: dX=%f, dP=%f, dS=%f", dX, dP, dS)
	}

	// Update variables using Euler's method
	X += dX * dt
	P += dP * dt
	S += dS * dt
	V += dV * dt

	return X, P, S, V, nil
}
