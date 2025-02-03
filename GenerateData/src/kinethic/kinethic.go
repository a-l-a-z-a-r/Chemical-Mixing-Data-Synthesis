package kinethic

import (
	"fmt"
	"math"
)

func KineticModelStep(X, P, S, V, F, T, muRef, qpRef, qsRef, EaMu, EaQp, EaQs, Kis, Pix, Pmx, Inhib, Ksp, dt float64) (float64, float64, float64, float64, error) {
	// Constants
	R := 8.314     // Gas constant (J/mol/K)
	TRef := 298.15 // Reference temperature (K)

	if T <= 0 {
		return X, P, S, V, fmt.Errorf("temperature T must be greater than 0; got %f", T)
	}

	if F >= 2 {
		return X, P, S, V, fmt.Errorf("flow rate can't be larger than 2 Liters  %f", F)
	}

	Xmix := (F * X) / V // g/L Biomass concentreation(Even if a formula is given int Missing x-mixture data asumptions, paper does not provide sufficient data
	Pmix := (F * P) / V // g/L Lactic acid conentration ( formula not given obtained by mass balance)
	Smix := (F * S) / V // g/L Lacotose conentration ( fromula not given obtained by mass balance)

	// Temperature-dependent parameters using Arrhenius equation
	muMax := muRef * math.Exp(-EaMu/R*(1/T-1/TRef))
	qpMax := qpRef * math.Exp(-EaQp/R*(1/T-1/TRef))
	qsMax := qsRef * math.Exp(-EaQs/R*(1/T-1/TRef))

	// Differential equations
	dX := muMax*X*(1-(P-Pix)/(Pmx-Pix)) + F*X/V + Xmix
	dP := Inhib*dX + qpMax*X*S/(V*Ksp+S) + F*P/V - Inhib*(X/V) + Pmix
	dS := -qsMax*X*V*Kis/(V*Kis+S) + F*(S/V) + Smix
	dV := F

	// Update variables using Euler's method
	X += dX * dt
	P += dP * dt
	S += dS * dt
	V += dV * dt

	if X <= 0 {
		var X float64 = 0
		return X, P, S, V, nil
	}

	return X, P, S, V, nil
}
