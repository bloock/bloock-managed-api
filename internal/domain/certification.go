package domain

import "github.com/bloock/bloock-sdk-go/v2/entity/integrity"

type Certification struct {
	anchorID int
	anchor   *integrity.Anchor
	hash     string
	dataID   string
	data     []byte
}

func NewPendingCertification(data []byte) *Certification {
	return &Certification{data: data}
}

func NewCertification(anchorID int, hash string, anchor *integrity.Anchor) *Certification {
	return &Certification{anchorID: anchorID, hash: hash, anchor: anchor}
}

func (c Certification) AnchorID() int {
	return c.anchorID
}

func (c *Certification) SetAnchorID(anchorID int) {
	c.anchorID = anchorID
}

func (c Certification) Anchor() *integrity.Anchor {
	return c.anchor
}

func (c *Certification) SetAnchor(anchor *integrity.Anchor) {
	c.anchor = anchor
}

func (c Certification) Hash() string {
	return c.hash
}

func (c *Certification) SetHash(hash string) {
	c.hash = hash
}

func (c Certification) Data() []byte {
	return c.data
}

func (c *Certification) SetData(data []byte) {
	c.data = data
}

func (c Certification) DataID() string {
	return c.dataID
}

func (c *Certification) SetDataID(dataID string) {
	c.dataID = dataID
}
