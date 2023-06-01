package repository

import "github.com/bloock/bloock-sdk-go/v2/entity/integrity"

type NotificationRepository interface {
	NotifyCertification(hash string, anchor integrity.Anchor) error
}
