package dto

type UpdateJobStatusRequest struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
