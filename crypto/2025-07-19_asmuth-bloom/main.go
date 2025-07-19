package main

import (
	"fmt"
	"math/rand/v2"
	"sort"
)

func main() {
	secret := 42
	a := asmuthBloom{
		secretMaxRange: 43,
		threshold:      3,
		mods:           []int{101, 103, 107, 109, 113},
	}
	if err := runAsmuthBloom(a, secret); err != nil {
		panic(err)
	}
}

func runAsmuthBloom(a asmuthBloom, secret int) error {
	fmt.Printf("config: %+v\n", a)
	fmt.Println("secret:", secret)

	shares, err := a.generateShares(secret)
	if err != nil {
		return fmt.Errorf("generateShares: %w", err)
	}
	fmt.Printf("shares: %+v\n", shares)

	reconstructableParty := make([]share, a.threshold)
	for i := range a.threshold {
		reconstructableParty[i] = shares[i]
	}
	fmt.Println("reconstructableParty:", reconstructableParty)

	reconstructedSecret, err := a.reconstructSecret(reconstructableParty)
	if err != nil {
		return fmt.Errorf("reconstructSecret: %w", err)
	}
	fmt.Println("reconstructedSecret:", reconstructedSecret)

	unreconstructableParty := make([]share, a.threshold-1)
	for i := range unreconstructableParty {
		unreconstructableParty[i] = shares[len(shares)-1-i]
	}
	fmt.Println("unreconstructableParty:", unreconstructableParty)

	reconstructedSecret2, err := a.reconstructSecret(unreconstructableParty)
	if err != nil {
		return fmt.Errorf("reconstructSecret: %w", err)
	}
	fmt.Println("reconstructedSecret2:", reconstructedSecret2)

	return nil
}

type share struct {
	mod   int
	value int
}

type asmuthBloom struct {
	secretMaxRange int
	threshold      int
	mods           []int
}

func (a asmuthBloom) generateShares(secret int) ([]share, error) {
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

	if a.secretMaxRange < 1 || a.secretMaxRange >= a.mods[0] {
		return nil, fmt.Errorf("secretMaxRange is out of range")
	}

	if secret < 1 || secret >= a.secretMaxRange {
		return nil, fmt.Errorf("secret is out of range")
	}

	maxUnreconstructableParty := 1
	for i := range a.threshold - 1 {
		maxUnreconstructableParty *= a.mods[len(a.mods)-1-i]
	}
	minReconstructableParty := 1
	for i := range a.threshold {
		minReconstructableParty *= a.mods[i]
	}
	if a.secretMaxRange*maxUnreconstructableParty >= minReconstructableParty {
		return nil, fmt.Errorf("mods are not valid: %d * %d >= %d", a.secretMaxRange, maxUnreconstructableParty, minReconstructableParty)
	}

	maxRand := (minReconstructableParty - secret) / a.secretMaxRange
	r := rand.IntN(maxRand)
	y := secret + r*a.secretMaxRange

	shares := make([]share, len(a.mods))
	for i := range a.mods {
		shares[i] = share{
			mod:   a.mods[i],
			value: y % a.mods[i],
		}
	}

	return shares, nil
}

func (a asmuthBloom) reconstructSecret(shares []share) (int, error) {
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

	reconstuctedSecret %= a.secretMaxRange

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
