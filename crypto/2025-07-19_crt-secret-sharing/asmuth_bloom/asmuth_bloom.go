package asmuth_bloom

import (
	"crt_secret_sharing/crt"
	"fmt"
	"math/rand/v2"
	"sort"
)

type Share struct {
	mod   int
	value int
}

type Config struct {
	secretMod int
	threshold int
	mods      []int
}

func (a Config) generateShares(secret int) ([]Share, error) {
	if a.threshold < 1 || a.threshold > len(a.mods) {
		return nil, fmt.Errorf("threshold is out of range")
	}

	modsIsSorted := sort.SliceIsSorted(a.mods, func(i, j int) bool {
		return a.mods[i] < a.mods[j]
	})
	if !modsIsSorted {
		return nil, fmt.Errorf("mods are not sorted")
	}

	// modsが互いに素である
	for i := range a.mods {
		for j := i + 1; j < len(a.mods); j++ {
			if gcd, _, _ := crt.ExtendedEuclidean(a.mods[i], a.mods[j]); gcd != 1 {
				return nil, fmt.Errorf("mods are not coprime: %d and %d", a.mods[i], a.mods[j])
			}
		}
	}

	if a.secretMod < 1 || a.secretMod >= a.mods[0] {
		return nil, fmt.Errorf("secretMaxRange is out of range")
	}

	if secret < 1 || secret >= a.secretMod {
		return nil, fmt.Errorf("secret is out of range")
	}

	securityProduct := 1
	for i := range a.threshold - 1 {
		securityProduct *= a.mods[len(a.mods)-1-i]
	}
	encodedSecretUpperBound := 1
	for i := range a.threshold {
		encodedSecretUpperBound *= a.mods[i]
	}
	if a.secretMod*securityProduct >= encodedSecretUpperBound {
		return nil, fmt.Errorf("mods are not valid: %d * %d >= %d", a.secretMod, securityProduct, encodedSecretUpperBound)
	}

	maxRand := (encodedSecretUpperBound - secret) / a.secretMod
	r := rand.IntN(maxRand)
	encodedSecret := secret + r*a.secretMod

	shares := make([]Share, len(a.mods))
	for i := range a.mods {
		shares[i] = Share{
			mod:   a.mods[i],
			value: encodedSecret % a.mods[i],
		}
	}

	return shares, nil
}

func (a Config) reconstructSecret(shares []Share) (int, error) {
	// if len(shares) < a.threshold {
	// 	return 0, fmt.Errorf("shares are not enough")
	// }

	values := make([]int, len(shares))
	for i := range shares {
		values[i] = shares[i].value
	}
	mods := make([]int, len(shares))
	for i := range shares {
		mods[i] = shares[i].mod
	}

	reconstuctedSecret, err := crt.Solve(values, mods)
	if err != nil {
		return 0, fmt.Errorf("chineseRemainderTheorem: %w", err)
	}

	reconstuctedSecret %= a.secretMod

	return reconstuctedSecret, nil
}
