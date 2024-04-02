package domain

import (
	"encoding/hex"
	"errors"
	"github.com/bloock/bloock-managed-api/internal/platform/utils"
	"sort"
)

var ProofAssemblerError = errors.New("couldn't assemble proofs")

type Proof struct {
	Leaves []string
	Nodes  []string
	Bitmap string
	Depth  string
	Root   string
	Anchor interface{}
}

func (p Proof) AssembleProof(subProof []Proof) (Proof, error) {
	if len(subProof) == 0 {
		return p, nil
	}

	sort.SliceStable(subProof, func(i, j int) bool {
		return subProof[i].Root < subProof[j].Root
	})

	var resBitmap utils.Bitset

	bitmapBytes, err := hex.DecodeString(p.Bitmap)
	if err != nil {
		return Proof{}, ProofAssemblerError
	}

	bitmap := utils.BitsetFromBytes(bitmapBytes)

	depths := ""
	nodes := make([]string, 0)
	leaves := make([]string, 0)

	leafIndex := 0
	nodeIndex := 0
	depthIndex := 0
	bitmapIndex := 0
	substateIndex := 0

	for i := 0; i < len(p.Nodes)+len(p.Leaves); i++ {
		bit := bitmap.GetBit(i)

		if !bit {
			substateProof := subProof[substateIndex]

			if p.Leaves[leafIndex] == substateProof.Root {
				substateBitmapBytes, err := hex.DecodeString(substateProof.Bitmap)
				if err != nil {
					return Proof{}, ProofAssemblerError
				}
				substateBitmap := utils.BitsetFromBytes(substateBitmapBytes)

				for si, sv := range substateBitmap {
					if si < len(substateProof.Nodes)+len(substateProof.Leaves) {
						resBitmap.SetBit(bitmapIndex+si, sv)
					}
				}

				bitmapIndex = bitmapIndex + len(substateProof.Nodes) + len(substateProof.Leaves)

				nodes = append(nodes, substateProof.Nodes...)
				leaves = append(leaves, substateProof.Leaves...)

				// Calculate depth
				leafDepthBytes, err := hex.DecodeString(p.Depth[depthIndex*4+2 : depthIndex*4+4])
				if err != nil {
					return Proof{}, ProofAssemblerError
				}

				substateDepthBytes, err := hex.DecodeString(substateProof.Depth)
				if err != nil {
					return Proof{}, ProofAssemblerError
				}

				for i := 1; i < len(substateDepthBytes); i = i + 2 {
					resDepth := substateDepthBytes[i] + leafDepthBytes[0]
					if resDepth < leafDepthBytes[0] {
						depths = depths + hex.EncodeToString([]byte{0x01, resDepth})
					} else {
						depths = depths + hex.EncodeToString([]byte{0x00, resDepth})
					}
				}

				substateIndex = substateIndex + 1
				leafIndex = leafIndex + 1
				depthIndex = depthIndex + 1
			} else {
				resBitmap.SetBit(bitmapIndex, false)
				bitmapIndex = bitmapIndex + 1

				depths = depths + p.Depth[depthIndex*4:depthIndex*4+4]
				depthIndex = depthIndex + 1

				leaves = append(leaves, p.Leaves[leafIndex])
				leafIndex = leafIndex + 1
			}

		} else {

			resBitmap.SetBit(bitmapIndex, true)
			bitmapIndex = bitmapIndex + 1

			depths = depths + p.Depth[depthIndex*4:depthIndex*4+4]
			depthIndex = depthIndex + 1

			nodes = append(nodes, p.Nodes[nodeIndex])
			nodeIndex = nodeIndex + 1
		}

	}

	resBitmapBytes := resBitmap.ToBytes()

	return Proof{
		Leaves: leaves,
		Nodes:  nodes,
		Bitmap: hex.EncodeToString(resBitmapBytes),
		Depth:  depths,
		Root:   p.Root,
		Anchor: p.Anchor,
	}, nil
}
