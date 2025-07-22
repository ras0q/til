package asmuth_bloom

import (
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
			if gcd, _, _ := extendedEuclidean(a.mods[i], a.mods[j]); gcd != 1 {
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

	reconstuctedSecret, err := chineseRemainderTheorem(values, mods)
	if err != nil {
		return 0, fmt.Errorf("chineseRemainderTheorem: %w", err)
	}

	reconstuctedSecret %= a.secretMod

	return reconstuctedSecret, nil
}

// extendedEuclidean は拡張ユークリッド互除法を実装し、
// a*x + b*y = gcd(a, b) となる (gcd, x, y) を返す。
func extendedEuclidean(a, b int) (int, int, int) {
	if a == 0 {
		return b, 0, 1
	}
	gcd, x1, y1 := extendedEuclidean(b%a, a)
	x := y1 - (b/a)*x1
	y := x1
	return gcd, x, y
}

// modInverse は a (mod m) のモジュラ逆数を計算する。
func modInverse(a, m int) (int, error) {
	gcd, x, _ := extendedEuclidean(a, m)
	if gcd != 1 {
		return 0, fmt.Errorf("modular inverse does not exist")
	}
	// x が負の場合も考慮し、正の余りを返す
	return (x%m + m) % m, nil
}

// chineseRemainderTheorem は中国の剰余定理を解く。
// a: 余りのスライス (a_i)
// m: 法のスライス (m_i)
// x ≡ a_i (mod m_i) の解 x を返す。
func chineseRemainderTheorem(a, m []int) (int, error) {
	if len(a) != len(m) {
		return 0, fmt.Errorf("slices a and m must have the same length")
	}

	// すべての法 m_i の積 M を計算
	M := 1
	for _, mi := range m {
		M *= mi
	}

	var result int = 0
	for i := range a {
		Mi := M / m[i]
		invMi, err := modInverse(Mi, m[i])
		if err != nil {
			return 0, fmt.Errorf("modInverse: %w", err)
		}
		term := a[i] * Mi * invMi
		result += term
	}

	// 最小の非負の解を返す
	return result % M, nil
}
