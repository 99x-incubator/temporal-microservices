package dto

type DisableRequest struct {
	RobotID string `json:"robot_id"`
	UserID  string `json:"user_id"`
}
