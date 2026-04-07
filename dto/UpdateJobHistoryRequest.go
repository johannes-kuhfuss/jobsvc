package dto

type UpdateJobHistoryRequest struct {
	Message string `json:"message" san:"trim,xss"`
}
