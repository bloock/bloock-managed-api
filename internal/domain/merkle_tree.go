package domain

import (
	"encoding/hex"
	"fmt"
	"github.com/bloock/bloock-managed-api/internal/platform/utils"
	"math"
)

type MerkleTree struct {
	Root  string
	Proof map[string]MerkleTreeProof
}

type MerkleTreeProof struct {
	Siblings []string
	Path     uint32
}

type Stack struct {
	Sibling []byte
	Depth   int
}

func (m MerkleTreeProof) ConvertToBloockProof(hash string) (Proof, error) {
	var proof Proof

	rawSiblings := make([][]byte, len(m.Siblings))
	for _, sb := range m.Siblings {
		siblingBytes, err := hex.DecodeString(sb)
		if err != nil {
			return Proof{}, err
		}
		rawSiblings = append(rawSiblings, siblingBytes)
	}
	leave, err := hex.DecodeString(hash)
	if err != nil {
		return Proof{}, err
	}

	var depth []uint32
	var bitmap uint32
	stack := make([]Stack, 0)
	n := len(rawSiblings) - 1

	counter := n + 2 + 4 - ((n + 2) % 4) - 1
	for i := n; i >= 0; i-- {
		cc := uint32(math.Pow(2, float64(i)))
		if m.Path&cc == 0 {
			proof.Nodes = append(proof.Nodes, hex.EncodeToString(rawSiblings[i]))
			bitmap += uint32(math.Pow(2, float64(counter)))
			depth = append(depth, uint32(n-i+1))
			counter--
		} else {
			stack = append(stack, Stack{
				Sibling: rawSiblings[i],
				Depth:   n - i + 1,
			})
		}
	}

	proof.Leaves = append(proof.Leaves, hex.EncodeToString(leave))
	depth = append(depth, uint32(n+1))
	counter--

	for i := len(stack) - 1; i >= 0; i-- {
		proof.Nodes = append(proof.Nodes, hex.EncodeToString(stack[i].Sibling))
		bitmap += uint32(math.Pow(2, float64(counter)))
		depth = append(depth, uint32(stack[i].Depth))
		counter--
	}

	var depthU16 []uint16
	for _, p := range depth {
		if p > 65535 {
			// Handle the error case where the uint32 value is too large for uint16
			return Proof{}, fmt.Errorf("error: value too large for uint16: %d", p)
		}
		depthU16 = append(depthU16, uint16(p))
	}

	var depthU8 []byte
	for _, x := range depthU16 {
		depthU8 = append(depthU8, utils.Uint16ToBytes(x)...)
	}

	depthHex := hex.EncodeToString(depthU8)

	proof.Depth = depthHex
	preBitmap := fmt.Sprintf("%x", bitmap)
	if len(preBitmap)%2 != 0 {
		preBitmap = fmt.Sprintf("%s%s", preBitmap, "0")
	}
	proof.Bitmap = preBitmap

	return proof, nil
}
