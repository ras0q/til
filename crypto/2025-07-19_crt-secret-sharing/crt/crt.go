package crt

import "fmt"

// ExtendedEuclidean は拡張ユークリッド互除法を実装し、
// a*x + b*y = gcd(a, b) となる (gcd, x, y) を返す。
func ExtendedEuclidean(a, b int) (int, int, int) {
	if a == 0 {
		return b, 0, 1
	}
	gcd, x1, y1 := ExtendedEuclidean(b%a, a)
	x := y1 - (b/a)*x1
	y := x1
	return gcd, x, y
}

// modInverse は a (mod m) のモジュラ逆数を計算する。
func modInverse(a, m int) (int, error) {
	gcd, x, _ := ExtendedEuclidean(a, m)
	if gcd != 1 {
		return 0, fmt.Errorf("modular inverse does not exist")
	}
	// x が負の場合も考慮し、正の余りを返す
	return (x%m + m) % m, nil
}

// Solve は中国の剰余定理を解く。
// a: 余りのスライス (a_i)
// m: 法のスライス (m_i)
// x ≡ a_i (mod m_i) の解 x を返す。
func Solve(a, m []int) (int, error) {
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
