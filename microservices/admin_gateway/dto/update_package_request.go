package dto

type UpdatePackageRequest struct {
	UserID    string `json:"user_id"`
	PackageID string `json:"package_id"`
}
