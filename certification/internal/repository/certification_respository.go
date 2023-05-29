package repository

import "github.com/bloock/bloock-sdk-go/v2/entity/record"

type CertificationRepository interface {
	Certify(bytes []record.Record) (int, error)
}
