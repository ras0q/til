package mignotte

import (
	"crt_secret_sharing/crt"
	"fmt"
	"sort"
)

type Share struct {
	mod   int
	value int
}

type Config struct {
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

	allModsProduct := 1
	for i := range a.mods {
		allModsProduct *= a.mods[i]
	}
	if secret < 1 || secret >= allModsProduct {
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
	if securityProduct >= secret || secret >= encodedSecretUpperBound {
		return nil, fmt.Errorf("mods are not valid: %d >= %d or %d >= %d", securityProduct, secret, secret, encodedSecretUpperBound)
	}

	shares := make([]Share, len(a.mods))
	for i := range a.mods {
		shares[i] = Share{
			mod:   a.mods[i],
			value: secret % a.mods[i],
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

	allModsProduct := 1
	for i := range a.mods {
		allModsProduct *= a.mods[i]
	}

	reconstuctedSecret %= allModsProduct

	return reconstuctedSecret, nil
}
