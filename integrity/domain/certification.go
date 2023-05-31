package domain

import "github.com/bloock/bloock-sdk-go/v2/entity/integrity"

type Certification struct {
	anchorID int
	anchor   *integrity.Anchor
	hashes   []string
}

func NewCertification(anchorID int, hash []string, anchor *integrity.Anchor) *Certification {
	return &Certification{anchorID: anchorID, hashes: hash, anchor: anchor}
}

func (c Certification) AnchorID() int {
	return c.anchorID
}

func (c Certification) Anchor() *integrity.Anchor {
	return c.anchor
}

func (c Certification) Hashes() []string {
	return c.hashes
}
