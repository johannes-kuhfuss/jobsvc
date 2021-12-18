package dto

type UpdateJobStatusRequest struct {
	Status string `json:"status"`
	ErrMsg string `json:"err_msg"`
}
