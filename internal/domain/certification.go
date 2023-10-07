package domain

import "github.com/bloock/bloock-sdk-go/v2/entity/record"

type Certification struct {
	AnchorID int
	Hash     string
	DataID   string
	Data     []byte
	Record   *record.Record
}
