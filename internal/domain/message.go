package domain

type Message struct {
	Hash     string
	Root     string
	AnchorID int
	Proof    MerkleTreeProof
}
