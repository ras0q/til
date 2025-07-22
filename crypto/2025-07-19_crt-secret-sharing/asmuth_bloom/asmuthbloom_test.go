package asmuth_bloom

import (
	"testing"
)

// follow https://github.com/gotestyourself/gotest.tools/wiki/Go-Testing-Patterns
func Test_AsmuthBloom(t *testing.T) {
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
			t.Fatalf("reconstructSecret: %v", err)
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
				secretMod: 43,
				threshold: 3,
				mods:      []int{101, 103, 107, 109, 113},
			},
			secret: 42,
			chooseShares: func(shares []Share) []Share {
				return shares
			},
		},
		"has insufficient shares": {
			config: Config{
				secretMod: 43,
				threshold: 3,
				mods:      []int{101, 103, 107, 109, 113},
			},
			secret: 42,
			chooseShares: func(shares []Share) []Share {
				return shares[3:5]
			},
			shouldFail: true,
		},
		"has unknown shares": {
			config: Config{
				secretMod: 43,
				threshold: 3,
				mods:      []int{101, 103, 107, 109, 113},
			},
			secret: 42,
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
