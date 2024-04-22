package domain

import (
	"encoding/hex"
	"github.com/bloock/bloock-managed-api/internal/platform/utils"
)

type NormalizedProof struct {
	Depths utils.BitsetInt
	Bitmap utils.BitsetBool
	Hashes []string
	Leaves []string
	Root   string
	Anchor interface{}
}

func NormalizedProofFromBloockProof(proofs []BloockProof) ([]NormalizedProof, error) {
	normalizedProofs := make([]NormalizedProof, 0)

	for _, proof := range proofs {
		bitmapBytes, err := hex.DecodeString(proof.Bitmap)
		if err != nil {
			return nil, err
		}
		bitmap := utils.BitsetBoolFromBytes(bitmapBytes)

		depths, err := utils.BitsetIntFromString(proof.Depth)
		if err != nil {
			return nil, err
		}

		normalizedProofs = append(normalizedProofs, NormalizedProof{
			Depths: depths,
			Bitmap: bitmap,
			Hashes: proof.Nodes,
			Leaves: proof.Leaves,
			Root:   proof.Root,
			Anchor: proof.Anchor,
		})
	}

	return normalizedProofs, nil
}

func (p NormalizedProof) ToBloockProof() (BloockProof, error) {
	bitset := hex.EncodeToString(p.Bitmap.ToBytes())
	depths, err := p.Depths.ToString()
	if err != nil {
		return BloockProof{}, err
	}

	return BloockProof{
		Leaves: p.Leaves,
		Nodes:  p.Hashes,
		Bitmap: bitset,
		Depth:  depths,
		Root:   p.Root,
		Anchor: p.Anchor,
	}, nil
}
