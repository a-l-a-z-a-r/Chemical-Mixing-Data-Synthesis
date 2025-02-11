module github.com/main

go 1.23.5

replace github.com/kinethic => ../kinethic

replace github.com/simulation => ../simulation

replace github.com/phcalc => ../phcalc

require (
	github.com/phcalc v0.0.0-00010101000000-000000000000
	github.com/simulation v0.0.0-00010101000000-000000000000
)

require github.com/kinethic v0.0.0-00010101000000-000000000000 // indirect
