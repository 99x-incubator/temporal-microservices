package dto

// NotificationMessage represents the structure of the notification to be sent
type NotificationMessage struct {
	UserID  string `json:"user_id"`
	Message string `json:"message"`
}
