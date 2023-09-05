package domain

type Certification struct {
	AnchorID int
	Hash     string
	DataID   string
	Data     []byte
}
