// Ersoy, Oğuzhan, Thomas Brochmann Pedersen, and Emin Anarim. 2020.
// “Homomorphic Extensions of CRT-Based Secret Sharing.”
// Discrete Applied Mathematics (Amsterdam, Netherlands: 1988) 285 (October): 317–29.
// https://www.sciencedirect.com/science/article/pii/S0166218X20303012
package mignotte

import (
	"fmt"
	"math/rand/v2"
)

type HomomorphicV1Share struct {
	base       Share
	maskShares []Share
}

type HomomorphicV1Config struct {
	base             Config
	secrecyThreshold int
}

func (c HomomorphicV1Config) generateShares(secret int) ([]HomomorphicV1Share, error) {
	if c.base.threshold != len(c.base.mods) {
		return nil, fmt.Errorf("In homomorphic v1 extension, the threshold must be equal to the number of parties")
	}

	allModsProduct := 1
	for i := range c.base.mods {
		allModsProduct *= c.base.mods[i]
	}

	maskedSecret := secret
	masks := make([]int, c.secrecyThreshold)
	for i := range c.secrecyThreshold {
		masks[i] = rand.IntN(allModsProduct)
		maskedSecret = (maskedSecret + masks[i]) % allModsProduct
	}

	baseShares, err := c.base.generateShares(maskedSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to generate shares: %w", err)
	}

	shares := make([]HomomorphicV1Share, 0, len(baseShares))
	for i, share := range baseShares {
		rotatedMasks := make([]Share, c.secrecyThreshold)
		for j := range c.secrecyThreshold {
			mod := c.base.mods[(i+j)%len(c.base.mods)]
			rotatedMasks[j] = Share{
				mod:   mod,
				value: masks[j] % mod,
			}
		}

		shares = append(shares, HomomorphicV1Share{
			base:       share,
			maskShares: rotatedMasks,
		})
	}

	return shares, nil
}

func (c HomomorphicV1Config) reconstructSecret(shares []HomomorphicV1Share) (int, error) {
	baseShares := make([]Share, len(shares))
	maskShares := make([][]Share, c.secrecyThreshold)
	for i, share := range shares {
		baseShares[i] = share.base

		for j, maskShare := range share.maskShares {
			maskShares[j] = append(maskShares[j], maskShare)
		}
	}

	maskedSecret, err := c.base.reconstructSecret(baseShares)
	if err != nil {
		return 0, fmt.Errorf("failed to reconstruct secret: %w", err)
	}

	masks := make([]int, c.secrecyThreshold)
	for i := range c.secrecyThreshold {
		reconstructedMask, err := c.base.reconstructSecret(maskShares[i])
		if err != nil {
			return 0, fmt.Errorf("failed to reconstruct mask: %w", err)
		}
		masks[i] = reconstructedMask
	}

	allModsProduct := 1
	for i := range c.base.mods {
		allModsProduct *= c.base.mods[i]
	}
	reconstructedSecret := maskedSecret
	for i := range c.secrecyThreshold {
		reconstructedSecret = (reconstructedSecret - masks[i] + allModsProduct) % allModsProduct
	}

	return reconstructedSecret, nil
}
