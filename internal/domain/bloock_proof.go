package domain

import (
	"encoding/hex"
	"errors"
	"github.com/bloock/bloock-managed-api/internal/platform/utils"
	"math"
	"sort"
)

var ProofAssemblerError = errors.New("couldn't assemble proofs")

type BloockProof struct {
	Leaves []string
	Nodes  []string
	Bitmap string
	Depth  string
	Root   string
	Anchor interface{}
}

func (p BloockProof) AssembleProof(subProof []BloockProof) (BloockProof, error) {
	if len(subProof) == 0 {
		return p, nil
	}

	sort.SliceStable(subProof, func(i, j int) bool {
		return subProof[i].Root < subProof[j].Root
	})

	var resBitmap utils.BitsetBool

	bitmapBytes, err := hex.DecodeString(p.Bitmap)
	if err != nil {
		return BloockProof{}, ProofAssemblerError
	}

	bitmap := utils.BitsetBoolFromBytes(bitmapBytes)

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
					return BloockProof{}, ProofAssemblerError
				}
				substateBitmap := utils.BitsetBoolFromBytes(substateBitmapBytes)

				for si, sv := range substateBitmap {
					if si < len(substateProof.Nodes)+len(substateProof.Leaves) {
						resBitmap.SetBit(bitmapIndex+si, sv)
					}
				}

				bitmapIndex = bitmapIndex + len(substateProof.Nodes) + len(substateProof.Leaves)

				nodes = append(nodes, substateProof.Nodes...)
				leaves = append(leaves, substateProof.Leaves...)

				leafDepthBytes, err := hex.DecodeString(p.Depth[depthIndex*4+2 : depthIndex*4+4])
				if err != nil {
					return BloockProof{}, ProofAssemblerError
				}

				substateDepthBytes, err := hex.DecodeString(substateProof.Depth)
				if err != nil {
					return BloockProof{}, ProofAssemblerError
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

	return BloockProof{
		Leaves: leaves,
		Nodes:  nodes,
		Bitmap: hex.EncodeToString(resBitmapBytes),
		Depth:  depths,
		Root:   p.Root,
		Anchor: p.Anchor,
	}, nil
}

func JoinBloockMultiProofs(bloockProofs []BloockProof) (BloockProof, error) {
	proofs, err := NormalizedProofFromBloockProof(bloockProofs)
	if err != nil {
		return BloockProof{}, err
	}

	maxDepth := getMaxDepth(proofs[0].Depths)
	result := NormalizedProof{Root: proofs[0].Root}
	proofIterators := make([][3]int, len(proofs))
	var lastDepth *int

	finalizedProofs := 0
	j := 0

	for finalizedProofs < len(proofs) {

		var hashValue string
		result.Bitmap = append(result.Bitmap, true)
		result.Depths = append(result.Depths, 0)

		for i := 0; i < len(proofs); i++ {
			if proofIterators[i][0] != len(proofs[i].Depths) {
				if lastDepth != nil {
					if proofIterators[i][1] > 0 {
						proofIterators[i][1] -= int(math.Pow(2, float64(maxDepth-*lastDepth)))
					} else if proofs[i].Depths[proofIterators[i][0]] < *lastDepth {
						proofIterators[i][1] = int(math.Pow(2, float64(maxDepth-proofs[i].Depths[proofIterators[i][0]]))) - int(math.Pow(2, float64(maxDepth-*lastDepth)))
						proofIterators[i][0] += 1
					} else if proofs[i].Depths[proofIterators[i][0]] == *lastDepth {
						proofIterators[i][0] += 1
					}
				}

				if proofIterators[i][0] == len(proofs[i].Depths) {
					proofIterators[i][0] += 1
					finalizedProofs += 1
				} else if proofIterators[i][1] == 0 && proofIterators[i][0] < len(proofs[i].Depths) {
					if proofs[i].Depths[proofIterators[i][0]] > result.Depths[len(result.Depths)-1] {
						result.Depths[len(result.Depths)-1] = proofs[i].Depths[proofIterators[i][0]]
						result.Bitmap[len(result.Bitmap)-1] = proofs[i].Bitmap[proofIterators[i][0]]
						if result.Bitmap[len(result.Bitmap)-1] {
							hashValue = proofs[i].Hashes[proofIterators[i][0]-proofIterators[i][2]]
						} else {
							hashValue = proofs[i].Leaves[0]
							proofIterators[i][2] = 1
						}
					}
				}
			}
		}

		if result.Bitmap[len(result.Bitmap)-1] {
			result.Hashes = append(result.Hashes, hashValue)
		} else {
			result.Leaves = append(result.Leaves, hashValue)
		}

		if len(result.Depths) > 0 {
			lastDepth = &result.Depths[len(result.Depths)-1]
		}

		j++

	}

	result.Depths = result.Depths[:len(result.Depths)-1]
	result.Bitmap = result.Bitmap[:len(result.Bitmap)-1]
	result.Hashes = result.Hashes[:len(result.Hashes)-1]

	return result.ToBloockProof()
}

func getMaxDepth(depths []int) int {
	maxDepth := depths[0]
	for _, d := range depths {
		if d > maxDepth {
			maxDepth = d
		}
	}
	return maxDepth
}
