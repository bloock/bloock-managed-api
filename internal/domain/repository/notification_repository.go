package repository

type NotificationRepository interface {
	NotifyCertification(hash string, file []byte) error
}
