package mignotte

import (
	"testing"
)

func Test_Mignotte(t *testing.T) {
	type testCase struct {
		config       Config
		secret       int
		chooseShares func(shares []Share) []Share
		shouldFail   bool
	}

	run := func(t *testing.T, tc testCase) {
		shares, err := tc.config.generateShares(tc.secret)
		if err != nil {
			t.Fatalf("generateShares: %v", err)
		}
		t.Logf("shares: %+v", shares)

		chosenShares := tc.chooseShares(shares)
		reconstructedSecret, err := tc.config.reconstructSecret(chosenShares)
		if err != nil {
			if !tc.shouldFail {
				t.Fatalf("reconstructSecret: %v", err)
			}
			return
		}

		reconstructedSuccessfully := tc.secret == reconstructedSecret
		if tc.shouldFail && reconstructedSuccessfully {
			t.Fatalf("case should fail but reconstructed successfully (secret: %d, reconstructedSecret: %d)", tc.secret, reconstructedSecret)
		}
		if !tc.shouldFail && !reconstructedSuccessfully {
			t.Fatalf("case should not fail but reconstructed failed (secret: %d, reconstructedSecret: %d)", tc.secret, reconstructedSecret)
		}
	}

	testCases := map[string]testCase{
		"has enough shares": {
			config: Config{
				threshold: 3,
				mods:      []int{101, 103, 107, 109, 113},
			},
			secret: 13000, // 上からthreshold-1個のmodの積より大きい値でないといけない
			chooseShares: func(shares []Share) []Share {
				return shares[:3]
			},
		},
		"has insufficient shares": {
			config: Config{
				threshold: 3,
				mods:      []int{101, 103, 107, 109, 113},
			},
			secret: 13000,
			chooseShares: func(shares []Share) []Share {
				return shares[3:5]
			},
			shouldFail: true,
		},
		"has unknown shares": {
			config: Config{
				threshold: 3,
				mods:      []int{101, 103, 107, 109, 113},
			},
			secret: 13000,
			chooseShares: func(shares []Share) []Share {
				return []Share{
					shares[0],
					shares[1],
					{
						mod:   shares[2].mod,
						value: shares[2].value * 2,
					},
				}
			},
			shouldFail: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			run(t, tc)
		})
	}
}

// TODO: これも成功するの？？？
func Test_Mignotte_Add(t *testing.T) {
	type testCase struct {
		secret1 int
		secret2 int
	}

	mods := []int{101, 103, 107, 109, 113}
	threshold := 5
	config := Config{
		threshold: threshold,
		mods:      mods,
	}

	modProduct := 1
	for i := range mods {
		modProduct *= mods[i]
	}

	run := func(t *testing.T, tc testCase) {
		shares1, err := config.generateShares(tc.secret1)
		if err != nil {
			t.Fatalf("generateShares1: %v", err)
		}
		shares2, err := config.generateShares(tc.secret2)
		if err != nil {
			t.Fatalf("generateShares2: %v", err)
		}

		addedShares := make([]Share, threshold)
		for i := range threshold {
			if shares1[i].mod != shares2[i].mod {
				t.Fatalf("mod mismatch: %d != %d", shares1[i].mod, shares2[i].mod)
			}

			value := (shares1[i].value + shares2[i].value) % shares1[i].mod
			addedShares[i] = Share{
				mod:   shares1[i].mod,
				value: value,
			}
		}

		reconstructed, err := config.reconstructSecret(addedShares)
		if err != nil {
			t.Fatalf("reconstructSecret (homomorphic addition): %v", err)
		}
		want := (tc.secret1 + tc.secret2) % modProduct
		if reconstructed != want {
			t.Fatalf("homomorphic addition failed: want %d, got %d", want, reconstructed)
		}
	}

	testCases := map[string]testCase{
		"small values": {
			secret1: 13000,
			secret2: 14000,
		},
		"zero and value": {
			secret1: 0,
			secret2: 15000,
		},
		"max and min": {
			secret1: 13000,
			secret2: 1,
		},
		"both zero": {
			secret1: 0,
			secret2: 0,
		},
		"large values": {
			secret1: 1100000,
			secret2: 1200000,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			run(t, tc)
		})
	}
}

// TODO: これも成功するの？？？
func Test_Mignotte_AddShares(t *testing.T) {
	mods := []int{101, 103, 107, 109, 113}
	threshold := 5
	config := Config{
		threshold: threshold,
		mods:      mods,
	}

	type testCase struct {
		secret    int
		operation func(int) int
	}

	run := func(t *testing.T, tc testCase) {
		shares, err := config.generateShares(tc.secret)
		if err != nil {
			t.Fatalf("generateShares: %v", err)
		}

		// 各Shareに同じ演算を適用
		operatedShares := make([]Share, threshold)
		for i := range threshold {
			operatedShares[i] = Share{
				mod:   shares[i].mod,
				value: tc.operation(shares[i].value),
			}
		}

		reconstructed, err := config.reconstructSecret(operatedShares)
		if err != nil {
			t.Fatalf("reconstructSecret (add shares): %v", err)
		}
		modProduct := 1
		for i := range mods {
			modProduct *= mods[i]
		}
		want := tc.operation(tc.secret) % modProduct
		if reconstructed != want {
			t.Fatalf("add shares failed: want %d, got %d", want, reconstructed)
		}
	}

	testCases := map[string]testCase{
		"add 1":     {secret: 13000, operation: func(x int) int { return x + 1 }},
		"add 100":   {secret: 13000, operation: func(x int) int { return x + 100 }},
		"add 0":     {secret: 13000, operation: func(x int) int { return x + 0 }},
		"add large": {secret: 13000, operation: func(x int) int { return x + 10000 }},
		// "mul 2":     {secret: 13000, operation: func(x int) int { return x * 2 }},
		// "mul 100":   {secret: 13000, operation: func(x int) int { return x * 100 }},
		// "mul 0":     {secret: 13000, operation: func(x int) int { return x * 0 }},
		// "mul large": {secret: 13000, operation: func(x int) int { return x * 10000 }},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			run(t, tc)
		})
	}
}
