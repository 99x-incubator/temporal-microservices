package dto

type PaymentConfirmationRequest struct {
	UserID    string `json:"user_id"`
	PackageID string `json:"package_id"`
}
