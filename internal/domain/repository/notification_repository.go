package repository

type NotificationRepository interface {
	NotifyCertification(hash string, whResponse interface{}, file []byte) error
}
