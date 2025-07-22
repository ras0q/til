package mignotte

import (
	"testing"
)

func Test_HomomorphicV1(t *testing.T) {
	type testCase struct {
		config       HomomorphicV1Config
		secret       int
		chooseShares func(shares []HomomorphicV1Share) []HomomorphicV1Share
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
			config: HomomorphicV1Config{
				base: Config{
					threshold: 5,
					mods:      []int{101, 103, 107, 109, 113},
				},
				secrecyThreshold: 3,
			},
			secret: 13000,
			chooseShares: func(shares []HomomorphicV1Share) []HomomorphicV1Share {
				return shares[:5]
			},
		},
		"has insufficient shares": {
			config: HomomorphicV1Config{
				base: Config{
					threshold: 5,
					mods:      []int{101, 103, 107, 109, 113},
				},
				secrecyThreshold: 3,
			},
			secret: 13000,
			chooseShares: func(shares []HomomorphicV1Share) []HomomorphicV1Share {
				return shares[3:5]
			},
			shouldFail: true,
		},
		"has unknown shares": {
			config: HomomorphicV1Config{
				base: Config{
					threshold: 5,
					mods:      []int{101, 103, 107, 109, 113},
				},
				secrecyThreshold: 3,
			},
			secret: 13000,
			chooseShares: func(shares []HomomorphicV1Share) []HomomorphicV1Share {
				corrupted := make([]HomomorphicV1Share, 5)
				copy(corrupted, shares[:5])
				corrupted[2].base.value = corrupted[2].base.value * 2
				return corrupted
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

func Test_HomomorphicV1_Add(t *testing.T) {
	type testCase struct {
		secret1 int
		secret2 int
	}

	mods := []int{101, 103, 107, 109, 113}
	threshold := 5
	config := HomomorphicV1Config{
		base: Config{
			threshold: threshold,
			mods:      mods,
		},
		secrecyThreshold: threshold,
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

		addedShares := make([]HomomorphicV1Share, threshold)
		for i := range threshold {
			base := Share{
				mod:   shares1[i].base.mod,
				value: (shares1[i].base.value + shares2[i].base.value) % shares1[i].base.mod,
			}
			maskShares := make([]Share, threshold)
			for j := range threshold {
				maskShares[j] = Share{
					mod:   shares1[i].maskShares[j].mod,
					value: (shares1[i].maskShares[j].value + shares2[i].maskShares[j].value) % shares1[i].maskShares[j].mod,
				}
			}
			addedShares[i] = HomomorphicV1Share{
				base:       base,
				maskShares: maskShares,
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

func Test_HomomorphicV1_AddShares(t *testing.T) {
	mods := []int{101, 103, 107, 109, 113}
	threshold := 5
	config := HomomorphicV1Config{
		base: Config{
			threshold: threshold,
			mods:      mods,
		},
		secrecyThreshold: threshold,
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
		operatedShares := make([]HomomorphicV1Share, threshold)
		for i := range threshold {
			operatedMaskShares := make([]Share, threshold)
			for j := range operatedMaskShares {
				operatedMaskShares[j] = Share{
					mod:   shares[i].maskShares[j].mod,
					value: tc.operation(shares[i].maskShares[j].value) % shares[i].maskShares[j].mod,
				}
			}

			operatedShares[i] = HomomorphicV1Share{
				base: Share{
					mod:   shares[i].base.mod,
					value: tc.operation(shares[i].base.value) % shares[i].base.mod,
				},
				maskShares: operatedMaskShares,
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
