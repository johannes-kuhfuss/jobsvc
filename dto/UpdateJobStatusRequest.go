package dto

type UpdateJobStatusRequest struct {
	Status  string `json:"status" san:"trim,xss,lower"`
	Message string `json:"message" san:"trim,xss"`
}
