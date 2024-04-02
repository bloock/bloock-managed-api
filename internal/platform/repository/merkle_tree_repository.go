package repository

import (
	"context"
	"encoding/hex"
	"github.com/bloock/bloock-managed-api/internal/domain"
	"github.com/bloock/bloock-managed-api/internal/platform/utils"
	"github.com/rs/zerolog"
	mt "github.com/txaty/go-merkletree"
)

type MerkleTreeRepository struct {
	config *mt.Config
	logger zerolog.Logger
}

func NewMerkleTreeRepository(l zerolog.Logger) MerkleTreeRepository {
	logger := l.With().Caller().Str("component", "merkle-tree-repository").Logger()

	cnf := &mt.Config{
		Mode:               mt.ModeProofGenAndTreeBuild,
		DisableLeafHashing: true,
		HashFunc:           utils.MerkleHashFunc,
	}

	return MerkleTreeRepository{
		config: cnf,
		logger: logger,
	}
}

type mtData struct {
	data string
}

func (m *mtData) Serialize() ([]byte, error) {
	return hex.DecodeString(m.data)
}

func mapToMtData(messages []domain.Message) []mt.DataBlock {
	mtd := make([]mt.DataBlock, 0)

	for _, c := range messages {
		mtd = append(mtd, &mtData{
			data: c.Hash,
		})
	}

	return mtd
}

func mapToMerkleTree(root []byte, merkleProofs []*mt.Proof, messages []domain.Message) domain.MerkleTree {
	proofs := make(map[string]domain.MerkleTreeProof, 0)

	hexRoot := hex.EncodeToString(root)

	for i, prf := range merkleProofs {
		siblings := make([]string, len(prf.Siblings))
		for _, sib := range prf.Siblings {
			siblings = append(siblings, hex.EncodeToString(sib))
		}

		proof := domain.MerkleTreeProof{
			Siblings: siblings,
			Path:     prf.Path,
		}
		proofs[messages[i].Hash] = proof
	}

	return domain.MerkleTree{
		Root:  hexRoot,
		Proof: proofs,
	}
}

func (m MerkleTreeRepository) Create(ctx context.Context, messages []domain.Message) (domain.MerkleTree, error) {
	merkle, err := mt.New(m.config, mapToMtData(messages))
	if err != nil {
		m.logger.Error().Err(err).Msg("")
		return domain.MerkleTree{}, err
	}

	return mapToMerkleTree(merkle.Root, merkle.Proofs, messages), nil
}
