package phcalc

import (
	"math"
)

// HendersonHasselbalch calculates pH using the buffer equation
func HendersonHasselbalch(acidConcentration, baseConcentration, pKa float64) float64 {
	if acidConcentration <= 0 {
		return 7.0 // Neutral pH when no acid is present
	}
	hPlus := pKa + math.Log10(baseConcentration/acidConcentration)
	return -math.Log10(math.Pow(10, -hPlus))
}

// CalculatePH estimates the pH based on the buffer model
func CalculatePH(lacticAcid, aceticAcid, bufferCapacity float64) float64 {
	const pKaLactic = 3.86 // pKa of lactic acid
	const pKaAcetic = 4.76 // pKa of acetic acid

	if lacticAcid <= 0 && aceticAcid <= 0 {
		return 7.0 // Neutral pH when no acids are present
	}

	// Using the Henderson-Hasselbalch equation for a combined buffer system
	phLactic := HendersonHasselbalch(lacticAcid, bufferCapacity, pKaLactic)
	phAcetic := HendersonHasselbalch(aceticAcid, bufferCapacity, pKaAcetic)

	// Weighted estimation based on acid concentrations and buffer capacity
	totalAcid := lacticAcid + aceticAcid
	if totalAcid > 0 {
		return ((phLactic * lacticAcid) + (phAcetic * aceticAcid)) / totalAcid
	}
	return 7.0
}
