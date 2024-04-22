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

func (m MerkleTreeProof) ParseToBloockProof(hash string, root string) (BloockProof, error) {
	var proof BloockProof

	rawSiblings := make([][]byte, 0)
	for _, sb := range m.Siblings {
		siblingBytes, err := hex.DecodeString(sb)
		if err != nil {
			return BloockProof{}, err
		}
		rawSiblings = append(rawSiblings, siblingBytes)
	}
	leave, err := hex.DecodeString(hash)
	if err != nil {
		return BloockProof{}, err
	}

	var depth utils.BitsetInt
	var bitmap uint32
	stack := make([]Stack, 0)
	n := len(rawSiblings) - 1

	counter := n + 2 + 4 - ((n + 2) % 4) - 1
	for i := n; i >= 0; i-- {
		cc := uint32(math.Pow(2, float64(i)))
		if m.Path&cc == 0 {
			proof.Nodes = append(proof.Nodes, hex.EncodeToString(rawSiblings[i]))
			bitmap += uint32(math.Pow(2, float64(counter)))
			depth = append(depth, n-i+1)
			counter--
		} else {
			stack = append(stack, Stack{
				Sibling: rawSiblings[i],
				Depth:   n - i + 1,
			})
		}
	}

	proof.Leaves = append(proof.Leaves, hex.EncodeToString(leave))
	depth = append(depth, n+1)
	counter--

	for i := len(stack) - 1; i >= 0; i-- {
		proof.Nodes = append(proof.Nodes, hex.EncodeToString(stack[i].Sibling))
		bitmap += uint32(math.Pow(2, float64(counter)))
		depth = append(depth, stack[i].Depth)
		counter--
	}

	depthString, err := depth.ToString()
	if err != nil {
		return BloockProof{}, err
	}
	proof.Depth = depthString

	preBitmap := fmt.Sprintf("%x", bitmap)
	if len(preBitmap)%2 != 0 {
		preBitmap = fmt.Sprintf("%s%s", preBitmap, "0")
	}
	proof.Bitmap = preBitmap
	proof.Root = root

	return proof, nil
}
