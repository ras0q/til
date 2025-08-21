package beavertriples

import (
	"fmt"
	"math/rand/v2"
)

func main() {
	secret := uint(0)
	fmt.Println(secret, -secret)

	msb := beaverTriples(secret)
	fmt.Println("Most Significant Bit:", msb)

	msb2 := beaverTriples(-secret)
	fmt.Println("Most Significant Bit:", msb2)
}

func beaverTriples(secret uint) uint {
	x, y := generateAdditiveShares(secret)
	xs, xc := generateXorShares(x)
	ys, yc := generateXorShares(y)

	var carryS, carryC uint

	// Beaver Triples
	for i := range 64 {
		ts, tc := generateTripleShares()

		// server
		xis := (xs >> i) & 1
		yis := (ys >> i) & 1
		ds := xis ^ ts.a
		es := yis ^ ts.b

		// client
		xic := (xc >> i) & 1
		yic := (yc >> i) & 1
		dc := xic ^ tc.a
		ec := yic ^ tc.b

		// share [d] and [e]
		// each party
		d := ds ^ dc
		e := es ^ ec

		// calculate gi
		// gi := xi ^ yi

		// server
		gis := e*ts.a ^ d*ts.b ^ ts.c

		// client
		gic := e*tc.a ^ d*tc.b ^ tc.c ^ d*e

		// carry
		ts2, tc2 := generateTripleShares()

		// server
		pis := xis ^ yis
		ds2 := pis ^ ts2.a
		es2 := carryS ^ ts2.b

		// client
		pic := xic ^ yic
		dc2 := pic ^ tc2.a
		ec2 := carryC ^ tc2.b

		// share [d2] and [e2]
		// each party
		d2 := ds2 ^ dc2
		e2 := es2 ^ ec2

		// server
		gis2 := e2*ts2.a ^ d2*ts2.b ^ ts2.c
		carryS = gis ^ gis2

		// client
		gic2 := e2*tc2.a ^ d2*tc2.b ^ tc2.c ^ d2*e2
		carryC = gic ^ gic2
	}

	msbS := ((xs >> 63) & 1) ^ ((ys >> 63) & 1) ^ carryS
	msbC := ((xc >> 63) & 1) ^ ((yc >> 63) & 1) ^ carryC

	msb := msbS ^ msbC

	return msb
}

type triple struct {
	a, b, c uint
}

func generateTripleShares() (triple, triple) {
	a, b := rand.UintN(2), rand.UintN(2)
	c := a * b

	as, ac := generateXorBitShares(a)
	bs, bc := generateXorBitShares(b)
	cs, cc := generateXorBitShares(c)

	ts := triple{as, bs, cs}
	tc := triple{ac, bc, cc}

	return ts, tc
}

func generateAdditiveShares(secret uint) (uint, uint) {
	serverShare := rand.Uint()
	clientShare := secret - serverShare

	return serverShare, clientShare
}

func generateXorShares(secret uint) (uint, uint) {
	serverShare := rand.Uint()
	clientShare := secret ^ serverShare

	return serverShare, clientShare
}

func generateXorBitShares(secret uint) (uint, uint) {
	serverShare := rand.UintN(2)
	clientShare := secret ^ serverShare

	return serverShare, clientShare
}
