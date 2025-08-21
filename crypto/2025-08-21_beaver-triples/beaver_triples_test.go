package beavertriples

import (
	"testing"
)

func FuzzBeaverTriples(f *testing.F) {
	tc := []int{0, 1, 10, 100, -1, -10, -100, 123456789, -123456789}
	for _, i := range tc {
		f.Add(i)
	}

	f.Fuzz(func(t *testing.T, i int) {
		msb := beaverTriples(uint(i))
		expected := uint(i) >> 63
		if msb != expected {
			t.Fatalf("MSB(%d) should be %d, but %d", i, expected, msb)
		}
	})
}
